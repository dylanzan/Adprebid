/**
* @Author: Dylan
* @Date: 2020/7/1 16:46
 */
package config

type Config struct {
	Basic   `yaml:"Basic"`
	Tencent `yaml:"Tencent"`
	Yiche   `yaml:"Yiche"`
	IQiyi   `yaml:"IQiyi"`
	Toutiao `yaml:"Toutiao"`
	Youku   `yaml:"Youku"`
	TvMao   `yaml:"TvMao"`
}

type Basic struct {
	ListenPort string
}

type Yiche struct {
	UpstreamAddrs       []string
	DefaultUpstreamAddr string
	TimesBackToSource   int
}

type Tencent struct {
	UpstreamAddrs       []string
	DefaultUpstreamAddr string
	TimesBackToSource   int
}

type IQiyi struct {
	UpstreamAddrs       []string
	DefaultUpstreamAddr string
	TimesBackToSource   int
}

type Toutiao struct {
	UpstreamAddrs       []string
	DefaultUpstreamAddr string
	TimesBackToSource   int
}

type Youku struct {
	UpstreamAddrs       []string
	DefaultUpstreamAddr string
	TimesBackToSource   int
}

type TvMao struct {
	UpstreamAddrs       []string
	DefaultUpstreamAddr string
	TimesBackToSource   int
}
