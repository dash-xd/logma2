package app

type PublishRequest struct {
    Type            string      `json:"type"`
    SentTimeUtc     interface{} `json:"sentTimeUtc"`
    Message         interface{} `json:"message"`
    ParentNamespace string      `json:"parentNamespace"`
    ChildNamespace  string      `json:"childNamespace"`
    Channel         string      `json:"channel"`
}
