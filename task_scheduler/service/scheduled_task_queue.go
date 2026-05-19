package service

type ScheduledTaskQueue []ScheduledTask

func (q ScheduledTaskQueue) Len() int { return len(q) }

func (q ScheduledTaskQueue) Less(i, j int) bool {
	return q[i].LessThan(q[j])
}

func (q ScheduledTaskQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q *ScheduledTaskQueue) Push(task interface{}) {
	*q = append(*q, task.(ScheduledTask))
}

func (q *ScheduledTaskQueue) Pop() interface{} {
	n := q.Len()
	task := (*q)[n-1]
	(*q) = (*q)[:n-1]
	return task
}
