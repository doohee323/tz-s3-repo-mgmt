<template>
  <div class="container-fluid">
    <div class="row">
      <div class="col-lg-12">
        <confirmation-modal v-if="isAdmin" id="confirmDeletion" ref="confirmDeletion"
                            :confirm="confirm"
                            :entries="selectedS3RepoEntries" :param-value="m => m.name" :title="messages.confirmDeletion" param-name="name">
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
          <label class="col-lg-5" for="s3repo">S3 Repo</label>
          <select id="s3repo" v-model="s3repo" class="form-control" data-parsley-required-message="Either device or gateway is required." name="s3repo"
                  required @change="requestData">
            <option v-for="item in s3repos" :key="item.id" :selected="item === s3repo" :value="item">
              {{ item }}
            </option>
          </select>
        </div>
        <div class="col-lg-2 form-group">
        </div>
        <div class="col-lg-3">
          <button v-if="isAdmin" class="btn btn-primary float-right" @click="reloadPage">
            <span class="fas fa-rotate-right"></span>
            Retrieve
          </button>
          <router-link v-if="isAdmin" :to="{name: 'Upload Image', params: {s3repo: s3repo}}"
                       class="btn btn-primary float-right">
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
  props: [],
  data() {
    const awsRegions = getCache('config', 'awsRegions').split(',');
    let awsRegion = getCache('session', 'awsRegion');
    if (awsRegion === '') {
      awsRegion = awsRegions[0];
    }
    const s3Repos = getCache('config', 's3Repos').split(',');
    let s3Repo = getCache('session', 's3repo');
    if (s3Repo === '' || s3Repos.indexOf(s3Repo) === -1) {
      s3Repo = s3Repos[0];
    }
    setCache('session', 'referer', '/s3Repos');
    return {
      apiUrl: getCache('config', 'apiUrl'),
      isAdmin: getCache('session', 'admin') === 'true',
      awsRegion: awsRegion,
      awsRegions: awsRegions,
      s3repo: s3Repo,
      s3repos: s3Repos,
      s3RepoEntries: [],
      actionButtons: [],
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
      this.requestData(this);
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
          $this.selectedS3RepoEntries = [];
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