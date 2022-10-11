package main

type Account struct {
	Alias   []string      `yaml:"alias"`
	Host    string        `yaml:"host"`
	Webhook WebhookConfig `yaml:"webhook"`
	IsTest  bool          `yaml:"isTest"`
}

type WebhookConfig struct {
	Id   string `yaml:"id"`
	Host string `yaml:"host"`
}
