package models

import (
	"fmt"
	"time"
)

type App struct {
	Name       string `json:"name"`
	GitAddress string `json:"git_address"`
}

type MultibranchWebhookTrigger struct {
	Status string `json:"status"`
	Data   struct {
		TriggerResults map[string]interface{}
	} `json:"data"`
}

type TriggerResultsData struct {
	Triggered bool   `json:"triggered"`
	Id        int64  `json:"id"`
	Url       string `json:"url"`
}

func JenkinsJobCfgFile(appName string, gitAddress string) string {
	return fmt.Sprintf(`<?xml version='1.1' encoding='UTF-8'?>
<org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject plugin="workflow-multibranch@2.21">
  <actions/>
  <description></description>
  <properties>
    <org.csanchez.jenkins.plugins.kubernetes.KubernetesFolderProperty plugin="kubernetes@1.15.4">
      <permittedClouds/>
    </org.csanchez.jenkins.plugins.kubernetes.KubernetesFolderProperty>
    <org.jenkinsci.plugins.pipeline.modeldefinition.config.FolderConfig plugin="pipeline-model-definition@1.3.8">
      <dockerLabel></dockerLabel>
      <registry plugin="docker-commons@1.14"/>
    </org.jenkinsci.plugins.pipeline.modeldefinition.config.FolderConfig>
  </properties>
  <folderViews class="jenkins.branch.MultiBranchProjectViewHolder" plugin="branch-api@2.4.0">
    <owner class="org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject" reference="../.."/>
  </folderViews>
  <healthMetrics>
    <com.cloudbees.hudson.plugins.folder.health.WorstChildHealthMetric plugin="cloudbees-folder@6.8">
      <nonRecursive>false</nonRecursive>
    </com.cloudbees.hudson.plugins.folder.health.WorstChildHealthMetric>
  </healthMetrics>
  <icon class="jenkins.branch.MetadataActionFolderIcon" plugin="branch-api@2.4.0">
    <owner class="org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject" reference="../.."/>
  </icon>
  <orphanedItemStrategy class="com.cloudbees.hudson.plugins.folder.computed.DefaultOrphanedItemStrategy" plugin="cloudbees-folder@6.8">
    <pruneDeadBranches>true</pruneDeadBranches>
    <daysToKeep>7</daysToKeep>
    <numToKeep>10</numToKeep>
  </orphanedItemStrategy>
  <triggers>
    <com.cloudbees.hudson.plugins.folder.computed.PeriodicFolderTrigger plugin="cloudbees-folder@6.8">
      <spec>H/15 * * * *</spec>
      <interval>3600000</interval>
    </com.cloudbees.hudson.plugins.folder.computed.PeriodicFolderTrigger>
    <com.igalg.jenkins.plugins.mswt.trigger.ComputedFolderWebHookTrigger plugin="multibranch-scan-webhook-trigger@1.0.1">
      <spec></spec>
      <token>%s</token>
    </com.igalg.jenkins.plugins.mswt.trigger.ComputedFolderWebHookTrigger>
  </triggers>
  <disabled>false</disabled>
  <sources class="jenkins.branch.MultiBranchProject$BranchSourceList" plugin="branch-api@2.4.0">
    <data>
      <jenkins.branch.BranchSource>
        <source class="jenkins.plugins.git.GitSCMSource" plugin="git@3.10.0">
          <id>%s-%d</id>
          <remote>%s</remote>
          <credentialsId>devops-use</credentialsId>
          <traits>
            <jenkins.plugins.git.traits.BranchDiscoveryTrait/>
            <jenkins.plugins.git.traits.SubmoduleOptionTrait>
              <extension class="hudson.plugins.git.extensions.impl.SubmoduleOption">
                <disableSubmodules>false</disableSubmodules>
                <recursiveSubmodules>true</recursiveSubmodules>
                <trackingSubmodules>false</trackingSubmodules>
                <reference></reference>
                <parentCredentials>true</parentCredentials>
              </extension>
            </jenkins.plugins.git.traits.SubmoduleOptionTrait>
          </traits>
        </source>
        <strategy class="jenkins.branch.DefaultBranchPropertyStrategy">
          <properties class="empty-list"/>
        </strategy>
      </jenkins.branch.BranchSource>
    </data>
    <owner class="org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject" reference="../.."/>
  </sources>
  <factory class="org.jenkinsci.plugins.workflow.multibranch.WorkflowBranchProjectFactory">
    <owner class="org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject" reference="../.."/>
    <scriptPath>Jenkinsfile</scriptPath>
  </factory>
</org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject>
`, appName, appName, time.Now().UnixNano(), gitAddress)
}
