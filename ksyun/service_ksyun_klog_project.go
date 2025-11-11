package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	klog "github.com/kingsoftcloud/sdk-go/v2/ksyun/client/klog/v20200731"
)

type KlogProjectService struct {
	client *KsyunClient
}

func (lg *KlogProjectService) ReadAndSetProjects(d *schema.ResourceData, r *schema.Resource) (err error) {
	var (
		listProjectsResp *klog.ListProjectsResponse
	)

	conn := lg.client.klogconn

	req := klog.NewListProjectsRequest()

	if page, ok := d.GetOk("page"); ok {
		p := page.(int)
		req.Page = &p
	}

	if size, ok := d.GetOk("size"); ok {
		s := size.(int)
		req.Size = &s
	}

	if name, ok := d.GetOk("project_name"); ok {
		n := name.(string)
		req.ProjectName = &n
	}

	if desc, ok := d.GetOk("description"); ok {
		t := desc.(string)
		req.Description = &t
	}

	listProjectsResp, err = conn.ListProjectsSend(req)
	if err != nil {
		return
	}

	// 处理项目列表数据
	projects := make([]interface{}, 0, len(listProjectsResp.Projects))
	for _, project := range listProjectsResp.Projects {
		proj := map[string]interface{}{
			"project_name":     *project.ProjectName,
			"iam_project_id":   *project.IamProjectId,
			"iam_project_name": *project.IamProjectName,
			"region":           *project.Region,
			"create_time":      *project.CreateTime,
			"update_time":      *project.UpdateTime,
			"status":           *project.Status,
			"log_pool_num":     *project.LogPoolNum,
		}

		// 处理标签
		tags := make([]map[string]interface{}, 0, len(project.Tags))
		for _, tag := range project.Tags {
			t := map[string]interface{}{
				"key":   *tag.Key,
				"value": *tag.Value,
			}
			tags = append(tags, t)
		}
		proj["tags"] = tags

		projects = append(projects, proj)
	}

	// 设置项目列表到ResourceData
	//if err := d.Set("projects", projects); err != nil {
	//	return fmt.Errorf("set projects failed: %s", err)
	//}
	//
	//d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  projects,
		idFiled:     "project_name",
		nameField:   "project_name",
		targetField: "projects",
		//extra: map[string]SdkResponseMapping{
		//	"ProjectName": {
		//		Field:    "project_name",
		//		KeepAuto: true,
		//	},
		//},
	})
}
