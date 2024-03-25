<template>
  <div class="container-fluid">
    <div class="row">
      <div class="col-lg-12">
        <confirmation-modal ref="confirmDeletion" id="confirmDeletion" :title="messages.confirmDeletion"
                            :entries="selectedS3RepoEntries"
                            param-name="name" :param-value="m => m.name" v-if="isAdmin" :confirm="confirm">
          <template #confirm>Delete</template>
        </confirmation-modal>
        <h2 class="page-header">
          {{ $route.name }}
        </h2>
      </div>
    </div>
    <div class="row">
      <div class="col-lg-6 offset-lg-1">
        <div class="alert ajax-alert alert-dismissable d-none">
          <button type="button" class="close" data-dismiss="alert" aria-hidden="true">&times;</button>
          <p id='ajax-msg'></p>
        </div>
      </div>
    </div>
    <div class="row form-main">
      <div class="col-lg-12 form-group form-inline s3repo">
        <div class="col-lg-3 form-group">
          <label for="awsRegion" class="col-lg-5">Region</label>
          <select name="awsRegion" id="awsRegion" class="form-control" v-model="awsRegion" @change="requestData"
                  required data-parsley-required-message="Either device or gateway is required.">
            <option v-for="item in awsRegions" :value="item" :key="item.id" :selected="item === awsRegion">
              {{ item }}
            </option>
          </select>
        </div>
        <div class="col-lg-4 form-group">
          <label for="s3repo" class="col-lg-5">S3 Repo</label>
          <select name="s3repo" id="s3repo" class="form-control" v-model="s3repo" @change="requestData"
                  required data-parsley-required-message="Either device or gateway is required.">
            <option v-for="item in s3repos" :value="item" :key="item.id" :selected="item === s3repo">
              {{ item }}
            </option>
          </select>
        </div>
        <div class="col-lg-2 form-group">
        </div>
        <div class="col-lg-3">
          <button class="btn btn-primary float-right" @click="reloadPage" v-if="isAdmin">
            <span class="fas fa-plus"></span>
            Retrieve
          </button>
          <router-link :to="{name: 'Upload Image', params: {s3repo: s3repo}}" class="btn btn-primary float-right" v-if="isAdmin">
            <span class="fas fa-plus"></span>
            Upload a File
          </router-link>
        </div>
      </div>
    </div>
    <div class="row">
      <div class="col-lg-12">
        <data-table ref="s3RepoDataTable" :entries="s3RepoEntries" :cells-data="cellsData" :page-lengths="[25,50,100]"
                    :condensed="true"
                    :striped="true" :action-buttons="actionButtons" :zero-records-message="messages.zeroRecords"
                    @click="clickRow" @loadeddata="mountedData" :overflow-wrap="true" :form-inline="false"
                    :columnDefs="repoTableColumnDefs"
                    @select="selectedS3RepoEntries=$event" @delete="$refs.confirmDeletion.show()">
          <tr>
            <th>CreateAt</th>
            <th data-class-name="enabled">FileType</th>
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
  props: [
  ],
  data() {
    const awsRegions = getCache('config', 'awsRegions').split(',');
    let awsRegion = getCache('session', 'awsRegion');
    if (awsRegion === '') {
      awsRegion = awsRegions[0];
    }
    const s3Repos = getCache('config', 's3Repos').split(',');
    let s3Repo = getCache('session', 's3repo');
    if (s3Repo === '') {
      s3Repo = s3Repos[0];
    }
    return {
      apiUrl: getCache('config', 'apiUrl'),
      isAdmin: getCache('session', 'admin') === 'true',
      awsRegion: awsRegion,
      awsRegions: awsRegions,
      s3repo: s3Repo,
      s3repos: s3Repos,
      s3RepoEntries: [],
      actionButtons: this.isAdmin ? [
        {
          label: 'Delete',
          icon: 'minus-square',
          context: 'danger',
          event: 'delete',
        },
      ] : null,
      messages: {
        confirmDeletion: 'Confirm Item Deletion',
        zeroRecords: 'No matching items found',
      },
      selectedS3RepoEntries: [],
      repoTableColumnDefs: [
        {"width": "20%", "targets": 0, "className": "text-center",},
        {"width": "10%", "targets": 1, "className": "text-center",},
        {"width": "40%", "targets": 2},
        {"width": "10%", "targets": 3, "className": "text-right",},
        {"width": "20%", "targets": 4, "className": "text-center",},
      ],
    }
  },
  methods: {
    reloadPage() {
      this.$router.go(0);
    },
    cellsData(e) {
      let link2 = `<a href="#" class="btn btn-primary btn-sm" @click="delete">
                                <span class="fas fa-trash"></span>Delete
                            </a>`;
      let downLink = `${this.apiUrl}/s3repo/download?s3repo=${this.s3repo}&object=${e.fileName}`
      return [
        {'data-order': e.updatedDt, 'html': e.updatedDt},
        {
          'data-order': e.fileType,
          'html': `<span class="fa-brands fa-${e.fileType}"></span>`
        },
        {'data-order': e.filePath, 'html': `<a href="${downLink}" target="_blank">${e.filePath}</a>`},
        parseFloat(e.size).toLocaleString('en-US'),
        this.isAdmin ? link2 : '',
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
      if ($this.$refs.s3RepoDataTable.dataTable) {
        $this.$refs.s3RepoDataTable.dataTable.clear();
        $this.$refs.s3RepoDataTable.dataTable.destroy();
      }
      $this.getItems((e) => {
        for (let i in e) {
          if (e[i].filePath.endsWith('apk')) {
            e[i]['fileType'] = 'android';
          } else if (e[i].filePath.endsWith('zip')) {
            e[i]['fileType'] = 'apple';
          } else {
            e[i]['fileType'] = 'question';
          }
        }
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
          $('.loading').hide();
          $this.$refs.confirmDeletion.hide();
        }
      });
    }
  },
  mounted() {
    checkLogin(this);
    this.requestData(this);
    initModal();
  },
};

</script>