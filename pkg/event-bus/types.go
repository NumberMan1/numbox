// Package eventbus defines portable topic and event-merge contracts.
// It does not implement publishing, subscription, or a request queue.
package eventbus

type Topic string

func (topic Topic) AddSuffix(suffix string) Topic { return topic + "_" + Topic(suffix) }
func (topic Topic) String() string                { return string(topic) }

// MergeableEventData may be merged only within the same topic and MergeKey.
// MergeKey should include the identity whose events may safely be combined.
type MergeableEventData interface {
	Merge(other MergeableEventData) MergeableEventData
	MergeKey() string
}
