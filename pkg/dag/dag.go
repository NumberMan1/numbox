package dag

import (
	"errors"
	"fmt"
)

// Node 表示 DAG 中的一个节点
type Node struct {
	ID       string           // 节点 ID
	Children map[string]*Node // 子节点的映射
}

// DAG 表示有向无环图
type DAG struct {
	Nodes map[string]*Node    // 图中所有节点的映射
	Edges map[string][]string // 边的映射，存储条件边
}

// NewDAG 创建一个新的空 DAG
func NewDAG() *DAG {
	return &DAG{
		Nodes: make(map[string]*Node),
		Edges: make(map[string][]string),
	}
}

// AddNode 向 DAG 中添加一个节点
func (dag *DAG) AddNode(id string) {
	if _, exists := dag.Nodes[id]; !exists {
		dag.Nodes[id] = &Node{
			ID:       id,
			Children: make(map[string]*Node),
		}
	}
}

// AddConditionalEdge 添加一条有条件的从 `from` 节点到 `to` 节点的有向边
// 只有当所有条件节点都存在时，目标节点才可达
func (dag *DAG) AddConditionalEdge(from, to string, conditions []string) error {
	if _, exists := dag.Nodes[from]; !exists {
		return errors.New("源节点不存在")
	}
	if _, exists := dag.Nodes[to]; !exists {
		return errors.New("目标节点不存在")
	}

	// 检查添加该边是否会导致环的产生
	if dag.createsCycle(from, to) {
		return errors.New("添加此边会导致环")
	}

	// 将条件边存储起来
	//key := fmt.Sprintf("%s->%s", from, to)
	key := from + "->" + to
	dag.Edges[key] = conditions

	// 检查所有条件节点是否存在
	for _, condition := range conditions {
		if _, exists := dag.Nodes[condition]; !exists {
			return nil // 如果条件节点不存在，不添加边
		}
	}

	// 如果所有条件节点都存在，则添加边
	dag.Nodes[from].Children[to] = dag.Nodes[to]
	return nil
}

// createsCycle 检查添加边是否会导致环，使用 DFS 实现
func (dag *DAG) createsCycle(from, to string) bool {
	visited := make(map[string]bool)  // 记录访问过的节点
	return dag.dfs(to, from, visited) // 深度优先搜索
}

// dfs 是一个辅助函数，用于在 createsCycle 中执行深度优先搜索
func (dag *DAG) dfs(current, target string, visited map[string]bool) bool {
	if current == target {
		return true // 找到环
	}
	if visited[current] {
		return false // 当前节点已访问过，结束递归
	}
	visited[current] = true // 标记当前节点为已访问

	// 遍历当前节点的子节点，递归检查是否存在环
	for childID := range dag.Nodes[current].Children {
		if dag.dfs(childID, target, visited) {
			return true
		}
	}
	return false // 未找到环
}

// GetDirectlyReachableNodes 返回从给定节点数组出发直接可到达的所有节点，不包括起始节点本身
func (dag *DAG) GetDirectlyReachableNodes(startNodes []string) ([]string, error) {
	// 如果 startNodes 为空，返回没有前置节点的节点
	if len(startNodes) == 0 {
		return dag.getNodesWithoutPredecessors(), nil
	}

	reachable := make(map[string]bool)

	for _, id := range startNodes {
		if _, exists := dag.Nodes[id]; !exists {
			return nil, fmt.Errorf("节点 %s 不存在", id)
		}
		for childID := range dag.Nodes[id].Children {
			if !contains(startNodes, childID) {
				// 检查条件边
				//key := fmt.Sprintf("%s->%s", id, childID)
				key := id + "->" + childID
				if conditions, exists := dag.Edges[key]; exists {
					if allConditionsMet(startNodes, conditions) {
						reachable[childID] = true
					}
				} else {
					reachable[childID] = true
				}
			}
		}
	}

	result := make([]string, 0, len(reachable))
	for node := range reachable {
		result = append(result, node)
	}

	return result, nil
}

// allConditionsMet 检查所有条件节点是否都在起始节点数组中
func allConditionsMet(startNodes, conditions []string) bool {
	for _, condition := range conditions {
		if !contains(startNodes, condition) {
			return false
		}
	}
	return true
}

// contains 检查切片中是否包含指定元素
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// getNodesWithoutPredecessors 返回没有前置节点的节点
func (dag *DAG) getNodesWithoutPredecessors() []string {
	inDegree := make(map[string]int)

	// 初始化所有节点的入度为0
	for id := range dag.Nodes {
		inDegree[id] = 0
	}

	// 计算节点的入度
	for _, node := range dag.Nodes {
		for childID := range node.Children {
			inDegree[childID]++
		}
	}

	// 找出入度为0的节点
	result := []string{}
	for id, degree := range inDegree {
		if degree == 0 {
			result = append(result, id)
		}
	}

	return result
}

// Print 打印 DAG 的结构
func (dag *DAG) Print() {
	visited := make(map[string]bool)
	for id := range dag.Nodes {
		if !visited[id] {
			dag.printNode(dag.Nodes[id], visited, 0, "")
		}
	}
}

// printNode 是一个递归辅助函数，用于打印节点及其子节点
func (dag *DAG) printNode(node *Node, visited map[string]bool, level int, prefix string) {
	if visited[node.ID] {
		return
	}
	visited[node.ID] = true

	// 打印当前节点，使用缩进表示层级
	if prefix == "" {
		fmt.Printf("%s%s\n", getIndent(level), node.ID)
	} else {
		fmt.Printf("%s%s->%s\n", getIndent(level), prefix, node.ID)
	}
	for _, child := range node.Children {
		dag.printNode(child, visited, level+1, node.ID)
	}
}

// getIndent 返回指定层级的缩进字符串
func getIndent(level int) string {
	indent := ""
	for i := 0; i < level; i++ {
		indent += "  "
	}
	return indent
}
