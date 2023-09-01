# 查询ALB
data "ksyun_albs" "default" {
  # 过滤条件：alb的ID
  ids = []
  # 过滤条件：vpcID
  vpc_id = []
  # 过滤条件：状态
  state = []
  output_file = "./temp/albs.json"
}

# 查询ALB的监听器
data "ksyun_alb_listeners" "ls" {
  # 过滤条件：监听器ID
  ids = []
  # 过滤条件：ALB的ID
  alb_id = []
  output_file = "./temp/alb_listeners.json"
}

# 查询ALB的转发策略
data "ksyun_alb_rule_groups" "default" {
  # 过滤条件：转发规则的ID
  ids = []
  # 过滤条件：监听器ID
  alb_listener_id=[]
  output_file="./temp/alb_rule_groups.json"
}

# 查询ALB的证书组
data "ksyun_alb_listener_cert_groups" "default" {
  # 过滤条件：证书组ID
  ids = []
  # 过滤条件：监听器ID
  alb_listener_id=[]
  output_file="./temp/alb_litener_cert_groups.json"
}




