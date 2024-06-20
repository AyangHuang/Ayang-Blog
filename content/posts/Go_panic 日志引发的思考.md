---
# 主页简介
# summary: ""
# 文章副标题
# subtitle: ""
# 作者信息
# author: ""
# authorLink: ""
# authorEmail: ""
# description: ""
# keywords: ""
# license: ""
# images: []
# 文章的特色图片
# featuredImage: ""
# 用在主页预览的文章特色图片
# featuredImagePreview: ""
# password:加密页面内容的密码，详见 主题文档 - 内容加密
# message:  加密提示信息，详见 主题文档 - 内容加密
linkToMarkdown: false
# 上面一般不用动
title: "panic 日志引发的思考"
date: 2024-06-01T01:12:45+08:00
lastmod: 2024-06-01T01:12:45+08:00
categories: ["Go"]
tags: []
---

最近在排查问题的时候，发现 panic 日志缺失等问题，所以就此问题展开如下探索

## 起因：测试环境 panic recover 日志平台找不到

最近在测试环境上线时，发现启动 panic 并无限循环

只显示重启次数和不健康的状态，要看 panic 日志的化得点击右侧的 STD 日志查看，就可以看到 panic 日志。随后并定位到代码问题

```go
func main() {
    ...
    loadConf()
    ...
    go cron_job()
}

func loadConf() {
    conf, err := getFromOss()
    if err != nil {
        // 定位到这里
        panic(fmt.Sprintf("s.reloadExemptConf ExemptReviewConfs conf load fail,id:%v,err:%+v", startId, e))
    }
}

func getFromOss() (*conf, err) {
	err := oss.Download()
	if err != nil {
		log.Errorc(ctx, "s.getExemptDetailFromBoss s.boss.Download confId(%+v) error(%+v)", conf.ID, err)
		return err
	}
    ...
}

// 使用的是 "github.com/robfig/cron" 包
func cron_job() {
    loadConf()
}
```

解释：oss 中的配置文件被删除了（定位不到为何删除，自动过期？人为删除？）导致的 painc，我重新上传后，实例就正常上线了

## 思考

1. 有 cron job 定时加载刷新配置，配置文件肯定是之前就丢了的。之所以之前不会 painc，肯定是 `github.com/robfig/cron` 在跑任务时有 recover painc

2. 但是为什么 recover panic 不会报警，问了 mentor，说测试环境 panic，recover panic 都不会报警（哈，实例挂掉也没人知道。。。）

3. 线上 painc，recover 报警是怎么发现 panic 的

## 探索

**思考一**

由于我在日志平台找不到关于 "panic" 关键词的日志，所以说 recover painc 的日志并没有被上报的日志平台。还好，在 `getFromOss()` 还有一个 error 日志

{{< image src="/images/panic 日志引发的思考/日志.png" width=100% height=100% caption="error 日志" >}} 

可以看到在 5.16 这天就已经获取不到日志，由于定时任务一直跑，所以很规律地打印日志

基本可以确定 panic 被 recover 了，看了下 `github.com/robfig/cron``github.com/robfig/cron` 的源码印证了猜想

{{< image src="/images/panic 日志引发的思考/日志打印.png" width=80% height=80% caption="recover 并打印日志" >}} 

而且由于 recover 后虽然打印了日志，但是只是输出到 STD，并不会上报到日志平台，也就查不到 panic 日志

**思考二和三**

原来使用公司工具包打印的日志，才会上报到公司的日志平台，**而 panic 日志使用内置的 print()，只能输入到 STD**，并不会上报

**那么公司的 painc 报警和 recover panic 报警是如何实现的呢**

其实很简单，由于 painc 日志只会在 STD 输出，并不会上报到平台。所以线上每一个容器内有一个进程，专门收集实例的 STD 日志，然后分析有没有 "panic" 字段并进行报警

**另外的思考** 

其实 cron job 的 recover painc 日志是有 STD 输出的，但是却不会触发报警。如果测试环境有 painc 日志报警，其实可以早点发现这个问题。所以我认为，测试环境也应该有 panic 日志的收集和报警

## 改进

`github.com/robfig/cron` 注入日志处理器，输入到 STD 且上报到日志平台。这样就不会存在找不到 cron job recover panic 但却在日志平台找不到日志的情况了

```go
import (
    goLog "log"
)

type ErrorLogger struct{}

func (ErrorLogger) Printf(format string, v ...interface{}) {
    // 公司的日志包，可以上报到公司的日志平台
	log.Error(format, v...)
    // go 内置的 log 包，默认是输出到 STD
    goLog.Printf(format, a...)
}

// 注入日志处理器
func (s *Service) loadproc(conf *configs.Config) { //nolint:unparam
	c := cron.New()
	c.ErrorLog = ErrorLogger{}
}
```

## 运用

最近处理一个脏数据的问题，运用到了这个知识点

场景：由于数据库表一个必填字段出现零值，mentor 让我排查所有该字段没有校验就写入的代码

我一看该字段被那么地方引用，头都大，后来转念一想，我只要在 dao 层加一个校验，如果为空就报警，并记录下堆栈，不就可以了么，这样就可以找到源头使得产生脏数据的原因，并且保证不会产生增量的脏数据。然后报警复用线上 painc 日志报警的能力就可以了，就不用自己写企微报警了（偷笑ing，我真是大聪明）

```go
var (
	errPTypeZero = errors.New("VideoTortRuleProperty.PType is empty")
)

// 打印含有"panic"日志到 std，目的是复用 std panic 日志报警能力，能及时知道并通过堆栈定位到导致 PType 产生脏数据的代码，杜绝脏数据入库
func logPanicPTypeZero(ctx context.Context) (err error) {
	buf := make([]byte, 64<<10)
	buf = buf[:runtime.Stack(buf, false)]
	fmt.Fprintf(os.Stderr, "panic recovered: err: %v\n%s", errPTypeZero, buf)
	log.Errorc(ctx, "panic recovered: err: %v\n%s", errPTypeZero, buf)
	// 向外抛出 err
	err = errPTypeZero
	return
}

func (d *Dao) SaveVideoTortRuleProperties(ctx context.Context, mm *dm.VideoTortRuleProperty) (err error) {
	// 判断是否脏数据
	if mm.PType.String() == dmMdl.PTNotMatch.String() {
		return logPanicPTypeZero(ctx)
	}
	orm := apmgorm.WithContext(ctx, d.videoOrm)
	return orm.Save(mm).Error
```

## 扩展

为什么 panic 时能打印日志，为了探索这个问题，我是用最常见的导致 panic 的方式：段错误。通过 debug 后我猜测，我发现其实段错误时会陷入操作系统内核态，并产生一个 syscall.SIGSEGV 标识为 11 的信号，而 Go 的 runtime 在启动时会注册所有信号的处理函数，然后转换为 Go 自己的信号，这样就能注册自己的处理函数。而 syscall.SIGSEGV 的默认处理函数是产生一个 panic，复用 panic 的机制来打印段错误的日志并结束进程

由于找不到任何资料有上面过程的解答，runtime 的源码也很难看。我就使用下面代码验证我的猜想

思路很简单，自己注册一个 syscall.SIGSEGV 的处理函数

```go
func main() {
	// 创建一个通道来接收信号
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	// 注册要接收的信号
	signal.Notify(sigs, syscall.SIGSEGV)

    sig := <-sigs
	fmt.Println("Received signal:", sig)
	fmt.Println("发生段错误，继续运行")
	for {
		time.Sleep(time.Second * 5)
		fmt.Println("发生段错误依然运行，不会产生 panic")
	}
}
```

```bash
go bulid -o test main.go
./test
Received signal: SIGSEGV
发生段错误，继续运行
发生段错误依然运行，不会产生 panic

# 另一个 bash
ps -ef | gerp test
ayang     215594  212789  0 23:35 pts/4    00:00:00 ./test
ayang     215640  149238  0 23:36 pts/0    00:00:00 grep --color=auto test
kill -11 215594  # 手动发送段错误信号
```
正当我欣喜若狂，觉得验证自己的猜想的时候。我又试了一下

```go
func main() {
	// 创建一个通道来接收信号
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	// 注册要接收的信号
	signal.Notify(sigs, syscall.SIGSEGV)

    go func() {
		time.Sleep(time.Second)
		var stu *Student
		stu.age = 18 // 代码触发段错误
	}()

    sig := <-sigs
	fmt.Println("Received signal:", sig)
	fmt.Println("发生段错误，继续运行")
	for {
		time.Sleep(time.Second * 5)
		fmt.Println("发生段错误依然运行，不会产生 panic")
	}
}
```

结果捕获不到段错误的信号

```bash
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x599260]
```

难道操作系统发出的段错误的信号和用户 kill 命令触发的段错误信号不同？然后 Go 根据是否是操作系统触发的来主动 painc，而用户注册的段错误信号处理器只针对用户主动产生的段错误信号？？？（有没有大佬有这方面知识的储备）

## 拓展思考

对 painc 有了更全面的思考，其实 panic 是 Go runtime 的一个机制，最后是主动打印日志并结束进程。注意：是主动的，可能是收到操作系统的段错误信号

当然如果如果是 9 信号，那么不会给程序机会，直接 kill 掉了，这是操作系统的机制。所以 painc 是 Go runtime 自己的机制，觉得是异常并影响到运行而主动结束进程，所以如果有 recover，也可以恢复正常。而 kill -9 则不会给程序任何机会

## End