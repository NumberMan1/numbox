# ioutil

WriteFull 补齐短写，拒绝零进展及非法字节计数，并返回 writer 的错误。
发生错误前已经写入的数据不会撤销；调用方负责 writer 的并发同步和关闭。
