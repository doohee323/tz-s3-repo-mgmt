<template>
  <div class="container-fluid">
    <div class="row">
      <div class="col-lg-12">
        <confirmation-modal id="confirmDeletion" ref="confirmDeletion" :confirm="confirm"
                            :entries="selectedS3RepoEntries"
                            :param-value="m => m.name" :title="messages.confirmDeletion" param-name="name">
          <template #confirm>Delete</template>
        </confirmation-modal>
        <confirmation-modal id="confirmRefresh" ref="confirmRefresh" :confirm="refresh"
                            :entries="selectedS3RepoEntries"
                            :param-value="m => m.name" :title="messages.confirmRefresh" param-name="name">
          <template #confirm>Refresh</template>
        </confirmation-modal>
        <h2 class="page-header">
          {{ $route.name }}
        </h2>
      </div>
    </div>
    <div class="row">
      <div class="col-lg-6 offset-lg-1">
        <div class="alert ajax-alert alert-dismissable d-none">
          <button aria-hidden="true" class="close" data-dismiss="alert" type="button">&times;</button>
          <p id='ajax-msg'></p>
        </div>
      </div>
    </div>
    <div class="row form-main">
      <div class="col-lg-12 form-group form-inline s3repo">
        <div class="col-lg-3 form-group">
          <label class="col-lg-5" for="awsRegion">Region</label>
          <select id="awsRegion" v-model="awsRegion" class="form-control" data-parsley-required-message="Either device or gateway is required." name="awsRegion"
                  required @change="requestData">
            <option v-for="item in awsRegions" :key="item.id" :selected="item === awsRegion" :value="item">
              {{ item }}
            </option>
          </select>
        </div>
        <div class="col-lg-4 form-group">
          <label class="col-lg-5" for="s3repo">S3 Study</label>
          <select id="s3repo" v-model="s3repo" class="form-control" data-parsley-required-message="Either device or gateway is required." name="s3repo"
                  required @change="requestData">
            <option v-for="item in s3repos" :key="item.id" :selected="item === s3repo" :value="item">
              {{ item }}
            </option>
          </select>
        </div>
        <div class="col-lg-1 form-group">
        </div>
        <div class="col-lg-4">
          <button class="btn btn-primary float-right" @click="refreshCache">
            <span class="fas fa-broom"></span>
            Refresh Cache
          </button>
          <button class="btn btn-primary float-right" @click="reloadPage">
            <span class="fas fa-rotate-right"></span>
            Retrieve
          </button>
          <router-link :to="{name: 'Upload Image', params: {s3repo: s3repo}}" class="btn btn-primary float-right">
            <span class="fas fa-upload"></span>
            Upload a File
          </router-link>
        </div>
      </div>
    </div>
    <div class="row">
      <div class="col-lg-12">
        <data-table ref="s3RepoDataTable" :cells-data="cellsData" :columnDefs="repoTableColumnDefs" :condensed="true"
                    :entries="s3RepoEntries"
                    :form-inline="false" :overflow-wrap="true"
                    :page-lengths="[25,50,100]" :striped="true" :zero-records-message="messages.zeroRecords" @click="clickRow"
                    @delete="$refs.confirmDeletion.show()"
                    @loadeddata="mountedData" @select="selectedS3RepoEntries=$event">
          <tr>
            <th>CreateAt</th>
            <th class="max-100r">FilePath</th>
            <th>Size</th>
            <th data-orderable="false">Control</th>
          </tr>
        </data-table>
      </div>
    </div>
  </div>
</template>

<script>

import $ from 'jquery';
import dataTable from '@/components/common/data-table.vue';
import confirmationModal from "@/components/common/confirmation-modal.vue";
import {getCache, setCache} from "@/assets/cache";
import {checkLogin, initModal, showMessage} from "@/assets/forms";

export default {
  components: {
    confirmationModal,
    dataTable,
  },
  props: [],
  data() {
    const awsRegions = getCache('config', 'awsRegions').split(',');
    let awsRegion = getCache('session', 'awsRegion');
    if (awsRegion === '') {
      awsRegion = awsRegions[0];
    }
    const s3Repos = getCache('config', 's3Tests').split(',');
    let s3Repo = getCache('session', 's3Repo');
    if (s3Repo === '' || s3Repos.indexOf(s3Repo) === -1) {
      s3Repo = s3Repos[0];
    }
    const cfTests = getCache('config', 'cfTests').split(',');
    let cloudfront = getCache('session', 'cloudfront');
    if (cloudfront === '' || cfTests.indexOf(cloudfront) === -1) {
      cloudfront = cfTests[0];
    }
    setCache('session', 'referer', '/test');
    return {
      apiUrl: getCache('config', 'apiUrl'),
      awsRegion: awsRegion,
      awsRegions: awsRegions,
      s3repo: s3Repo,
      s3repos: s3Repos,
      cloudfront: cloudfront,
      cfTests: cfTests,
      s3RepoEntries: [],
      actionButtons: [],
      messages: {
        confirmDeletion: 'Confirm Item Deletion',
        confirmRefresh: 'Confirm Cache Refresh',
        zeroRecords: 'No matching items found',
      },
      selectedS3RepoEntries: [],
      repoTableColumnDefs: [
        {"width": "20%", "targets": 0, "className": "text-center",},
        {"width": "40%", "targets": 1},
        {"width": "10%", "targets": 2, "className": "text-right",},
        {"width": "20%", "targets": 3, "className": "text-center",},
      ],
    }
  },
  methods: {
    reloadPage() {
      this.requestData(this);
    },
    refreshCache() {
      this.$refs.confirmRefresh.show();
    },
    cellsData(e) {
      let link2 = `<a href="#" class="btn btn-primary btn-sm" @click="delete">
                                <span class="fas fa-trash"></span>Delete
                            </a>`;
      let downLink = 'https://' + this.s3repo + e.fileName.substring(e.fileName.lastIndexOf('/') + 1, e.fileName.length);
      return [
        {'data-order': e.updatedDt, 'html': e.updatedDt},
        {'data-order': e.filePath, 'html': `<a href="${downLink}" target="_blank">${e.filePath}</a>`},
        parseFloat(e.size).toLocaleString('en-US'),
        link2,
      ];
    },
    clickRow(e, $currentRow) {
      const target = $(e.target);
      const item = this.s3RepoEntries[$currentRow];
      if (target.children().hasClass('fa-trash')) {
        this.deleteItem(item);
      }
    },
    mountedData() {
      // debugger;
    },
    getItems(cb) {
      let $this = this;
      $('.loading').show();
      const url = this.apiUrl + '/s3repos?awsRegion=' + this.awsRegion + '&s3repo=' + this.s3repo;
      $.ajax({
        url: url,
        type: "GET",
        success: function (res) {
          res = JSON.parse(res.message);
          $this.S3 = res;
          cb(res);
          $('.loading').hide();
        },
        error: function (res) {
          $('.loading').hide();
          res = JSON.parse(res.responseText);
          $("#ajax-msg").text(res.message);
          showMessage(true, 'warning');
        }
      });
    },
    requestData($this) {
      if ($this) $this = this;
      setCache('session', 'awsRegion', $this.awsRegion);
      setCache('session', 's3repo', $this.s3repo);
      setCache('session', 'cloudfront', $this.cloudfront);
      if ($this.$refs.s3RepoDataTable.dataTable) {
        $this.$refs.s3RepoDataTable.dataTable.clear();
        $this.$refs.s3RepoDataTable.dataTable.destroy();
      }
      $this.getItems((e) => {
        $this.s3RepoEntries = e;
        if ($this.$refs.s3RepoDataTable && $this.$refs.s3RepoDataTable.mPromise) {
          $this.$refs.s3RepoDataTable.mPromise.then(
              () => {
                $this.$refs.s3RepoDataTable.load();
              }
          );
        }
      });
    },
    deleteItem(item) {
      this.selectedS3RepoEntries.push(item);
      this.$refs.confirmDeletion.show();
    },
    confirm() {
      let $this = this;
      $('.loading').show();
      let apiUrl = getCache('config', 'apiUrl');
      const url = apiUrl + '/s3repo?s3repo=' + this.s3repo + '&object=' + this.selectedS3RepoEntries[0].fileName;
      $.ajax({
        url: url,
        type: "DELETE",
        success: function () {
          $this.requestData($this);
          $this.selectedS3RepoEntries = [];
          $('.loading').hide();
          $this.$refs.confirmDeletion.hide();
        }
      });
    },
    refresh() {
      let $this = this;
      $('.loading').show();
      debugger;
      let apiUrl = getCache('config', 'apiUrl');
      const url = apiUrl + '/s3repo/refreshCache?cloudfront=' + this.cloudfront + '&awsRegion=' + this.awsRegion;
      $.ajax({
        url: url,
        type: "POST",
        success: function () {
          $('.loading').hide();
          $this.$refs.confirmRefresh.hide();
        }
      });
    },
  },
  mounted() {
    checkLogin(this);
    this.requestData(this);
    initModal();
  },
};

</script>