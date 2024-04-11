<template>
  <div>
    <div class="pricing-header px-3 py-3 pt-md-5 pb-md-4 mx-auto text-center">
      <h1 class="display-5">CI/CD & Monitoring</h1>
      <p class="lead"></p>
    </div>
    <div class="container-fluid">
      <div class="col-lg-12">
        <div class="card card-default">
          <div class="card-body form-inline">
            <div class="col-sm-2"><h3>CI / CD flow</h3></div>
            <div class="col-sm-10 text-center">
              <img src="../assets/jenkins-argocd.png" usemap="#jenkins-argocd"/>
              <map name="jenkins-argocd">
                <area alt="Bitbucket" coords="125,0,180,100" href="https://github.com/doohee323/mirrortown.poi/admin/webhooks" shape="rect"
                      target="_blank">
                <area alt="Jenkins" coords="260,0,330,100" href="https://jenkins.default.eks-main-t.topzone.com/view/mtown/job/mirrortown.poi/job/develop/" shape="rect"
                      target="_blank">
                <area alt="ECR" coords="430,0,500,100" href="https://us-west-1.console.aws.amazon.com/ecr/repositories/private/336363860990/mtown-mirrortown-poi?region=us-west-1" shape="rect"
                      target="_blank">
                <area alt="Vault" coords="510,0,590,100" href="https://vault.default.eks-main-t.topzone.com/ui/vault/secrets/secret/show/mtown-dev/mirrortown.poi" shape="rect"
                      target="_blank">
                <area alt="Bitbucket" coords="260,145,330,250" href="https://github.com/doohee323/tz-argocd-repo/src/main/mtown-mirrortown-poi/develop/" shape="rect"
                      target="_blank">
                <area alt="ArgoCD" coords="430,145,500,250" href="https://argocd.default.eks-main-t.topzone.com/applications?search=poi&view=tiles&showFavorites=false&proj=mtown-dev&sync=&autoSync=&health=&namespace=&cluster=&labels=" shape="rect"
                      target="_blank">
                <area alt="EKS" coords="620,145,700,250" href="https://grafana.default.eks-main-t.topzone.com/d/6581e46e4e5c7ba40a07646395ef7b23/kubernetes-compute-resources-pod?orgId=1&refresh=10s&var-datasource=default&var-cluster=&var-namespace=mtown-dev&var-pod=mtown-mirrortown-poi-develop-79c44598b9-7j9rj" shape="rect"
                      target="_blank">
              </map>
            </div>
          </div>
        </div>
        <div class="card card-default">
          <div class="card-body form-inline">
            <div class="col-sm-2"><h3>Log / Monitor</h3></div>
            <div class="col-sm-10 text-center">
              <img src="../assets/prometheus-grafana.png" usemap="#prometheus-grafana"/>
              <map name="prometheus-grafana">
                <area alt="EKS" coords="5,30,70,100" href="https://grafana.default.eks-main-t.topzone.com/explore?orgId=1&left=%7B%22datasource%22:%22P8E80F9AEF21F6940%22,%22queries%22:%5B%7B%22refId%22:%22A%22,%22datasource%22:%7B%22type%22:%22loki%22,%22uid%22:%22P8E80F9AEF21F6940%22%7D,%22editorMode%22:%22builder%22,%22expr%22:%22%7Bapp%3D%5C%22mtown-mirrortown-poi-develop%5C%22%7D%20%7C%3D%20%60%60%22,%22queryType%22:%22range%22%7D%5D,%22range%22:%7B%22from%22:%22now-1h%22,%22to%22:%22now%22%7D%7D" shape="rect"
                      target="_blank">
                <area alt="Prometheus" coords="155,30,220,100" href="https://prometheus.default.eks-main-t.topzone.com/targets?search=mtown" shape="rect"
                      target="_blank">
                <area alt="AlertManager" coords="300,30,360,100" href="https://alertmanager.default.eks-main-t.topzone.com" shape="rect"
                      target="_blank">
                <area alt="Grafana" coords="155,160,220,230" href="https://grafana.default.eks-main-t.topzone.com/dashboards" shape="rect"
                      target="_blank">
              </map>
            </div>
          </div>
        </div>
        <div class="card card-default">
          <div class="card-body form-inline">
            <div class="col-sm-2"><h3>ELK / metrics</h3></div>
            <div class="col-sm-10 text-center">
              <img src="../assets/elk.png" usemap="#elk"/>
              <map name="elk">
                <area alt="EKS" coords="5,12,75,90" href="https://grafana.default.eks-main-t.topzone.com/explore?orgId=1&left=%7B%22datasource%22:%22P8E80F9AEF21F6940%22,%22queries%22:%5B%7B%22refId%22:%22A%22,%22datasource%22:%7B%22type%22:%22loki%22,%22uid%22:%22P8E80F9AEF21F6940%22%7D,%22editorMode%22:%22builder%22,%22expr%22:%22%7Bapp%3D%5C%22mtown-mirrortown-poi-develop%5C%22%7D%20%7C%3D%20%60%60%22,%22queryType%22:%22range%22%7D%5D,%22range%22:%7B%22from%22:%22now-1h%22,%22to%22:%22now%22%7D%7D" shape="rect"
                      target="_blank">
                <area alt="filebeat-1" coords="350,9,420,87" href="https://grafana.default.eks-main-t.topzone.com/explore?orgId=1&left=%7B%22datasource%22:%22P8E80F9AEF21F6940%22,%22queries%22:%5B%7B%22refId%22:%22A%22,%22datasource%22:%7B%22type%22:%22loki%22,%22uid%22:%22P8E80F9AEF21F6940%22%7D,%22editorMode%22:%22builder%22,%22expr%22:%22%7Bapp%3D%5C%22filebeat-filebeat%5C%22%7D%20%7C%3D%20%60%60%22,%22queryType%22:%22range%22%7D%5D,%22range%22:%7B%22from%22:%22now-1h%22,%22to%22:%22now%22%7D%7D" shape="rect"
                      target="_blank">
                <area alt="ALB" coords="6,114,76,192" href="https://us-west-1.console.aws.amazon.com/ec2/home?region=us-west-1#LoadBalancers:search=k8s-hypen-hypenhyp-947ca19836-1977768723.us-west-1.elb.amazonaws.com;sort=loadBalancerName" shape="rect"
                      target="_blank">
                <area alt="S3-1" coords="115,101,185,179" href="https://s3.console.aws.amazon.com/s3/buckets/alb-external-ingress-access-bucket-eks-main-t?region=us-west-1&tab=properties" shape="rect"
                      target="_blank">
                <area alt="SNS" coords="155,160,220,230" href="https://us-west-1.console.aws.amazon.com/sqs/v2/home?region=us-west-1#/queues/https%3A%2F%2Fsqs.us-west-1.amazonaws.com%2F336363860990%2Falb-external-ingress-access-bucket-event-queue" shape="rect"
                      target="_blank">
                <area alt="filebeat-2" coords="348,147,418,225" href="https://grafana.default.eks-main-t.topzone.com/explore?orgId=1&left=%7B%22datasource%22:%22P8E80F9AEF21F6940%22,%22queries%22:%5B%7B%22refId%22:%22A%22,%22datasource%22:%7B%22type%22:%22loki%22,%22uid%22:%22P8E80F9AEF21F6940%22%7D,%22editorMode%22:%22builder%22,%22expr%22:%22%7Bapp%3D%5C%22alb-external-ingress%5C%22%7D%20%7C%3D%20%60%60%22,%22queryType%22:%22range%22%7D%5D,%22range%22:%7B%22from%22:%22now-1h%22,%22to%22:%22now%22%7D%7D" shape="rect"
                      target="_blank">
                <area alt="SNS" coords="231,152,301,230" href="https://us-west-1.console.aws.amazon.com/sqs/v2/home?region=us-west-1#/queues/https%3A%2F%2Fsqs.us-west-1.amazonaws.com%2F336363860990%2Falb-external-ingress-access-bucket-event-queue" shape="rect"
                      target="_blank">
                <area alt="logstash" coords="454,6,524,84" href="https://logstash.elk.eks-main-t.topzone.com/" shape="rect"
                      target="_blank">
                <area alt="elasticsearch" coords="565,9,635,87" href="https://es.elk.eks-main-t.topzone.com/" shape="rect"
                      target="_blank">
                <area alt="kibana" coords="718,9,788,87" href="https://kibana.elk.eks-main-t.topzone.com/app/dashboards#/view/68d966d0-6347-11ee-bce1-9fbc05d24d84?_g=(filters:!(),refreshInterval:(pause:!t,value:0),time:(from:now-15m,to:now))&_a=(description:alb_dashboard_d0,filters:!(('$state':(store:appState),meta:(alias:!n,disabled:!f,index:'898b0570-5db4-11ee-bce1-9fbc05d24d84',key:client.address.keyword,negate:!t,params:(query:'43.201.181.72'),type:phrase),query:(match_phrase:(client.address.keyword:'43.201.181.72'))),('$state':(store:appState),meta:(alias:!n,disabled:!f,index:'898b0570-5db4-11ee-bce1-9fbc05d24d84',key:client.address.keyword,negate:!t,params:(query:'121.138.194.226'),type:phrase),query:(match_phrase:(client.address.keyword:'121.138.194.226'))),('$state':(store:appState),meta:(alias:!n,disabled:!f,index:'898b0570-5db4-11ee-bce1-9fbc05d24d84',key:client.address.keyword,negate:!t,params:(query:'39.7.25.235'),type:phrase),query:(match_phrase:(client.address.keyword:'39.7.25.235')))),fullScreenMode:!f,options:(hidePanelTitles:!f,useMargins:!t),query:(language:kuery,query:''),tags:!(),timeRestore:!f,title:alb_dashboard_d0,viewMode:view)" shape="rect"
                      target="_blank">
                <area alt="cloudfront" coords="7,209,77,287" href="https://us-east-1.console.aws.amazon.com/cloudfront/v4/home?region=us-west-1#/distributions/E3LZMG78XNBVP4/edit" shape="rect"
                      target="_blank">
                <area alt="S3-2" coords="115,203,185,281" href="https://s3.console.aws.amazon.com/s3/buckets/cloudfront-bucket-eks-main-t?region=us-west-1&tab=properties" shape="rect"
                      target="_blank">
              </map>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>

import $ from 'jquery';
import {checkLogin} from "@/assets/forms";
import {getCache} from "@/assets/cache";

export default {
  data() {
    this.$root.hasAuth = getCache('session', 'hasAuth') === 'true';
    return {
      hasAuth: this.$root.hasAuth
    };
  },
  props: [],
  mounted() {
    checkLogin(this);
    $('.loading').hide();
    $('.cluster-link').click(function (e) {
      let url = $(e.currentTarget).attr('href');
      url = url.replaceAll('tz-336363860990', $('#aws-account').val());
      window.open(url, "_blank");
    });
  },
};
</script>