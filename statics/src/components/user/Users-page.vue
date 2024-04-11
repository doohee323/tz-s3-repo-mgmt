<template>
  <div class="container-fluid">
    <div class="row">
      <div class="col-lg-12">
        <confirmation-modal v-if="isAdmin" id="confirmDeletion" ref="confirmDeletion"
                            :confirm="confirm"
                            :entries="selectedUserEntries" :param-value="m => m.name" :title="messages.confirmDeletion" param-name="name">
          <template #confirm>Delete</template>
        </confirmation-modal>
        <h2 class="page-header">
          {{ $route.name }}
        </h2>
      </div>
    </div>
    <div class="row form-main">
      <div class="col-lg-2 offset-lg-10 form-group form-inline s3repo">
        <router-link v-if="isAdmin" :to="{name: 'Add User'}" class="btn btn-primary float-right">
          <span class="fas fa-plus"></span>
          Add User
        </router-link>
        <button v-if="isAdmin" class="btn btn-primary float-right" @click="reloadPage">
          <span class="fas fa-plus"></span>
          Retrieve
        </button>
      </div>
    </div>
    <div class="row">
      <div class="col-lg-12">
        <data-table id="usersDataTable" ref="usersDataTable" :cells-data="cellsData" :columnDefs="repoTableColumnDefs"
                    :condensed="true" :entries="usersEntries"
                    :form-inline="false" :overflow-wrap="true"
                    :page-lengths="[25,50,100]" :striped="true" :zero-records-message="messages.zeroRecords" @click="clickRow"
                    @delete="$refs.confirmDeletion.show()"
                    @loadeddata="mountedData" @select="selectedUserEntries=$event">
          <tr>
            <th>Username</th>
            <th>Name</th>
            <th>Email</th>
            <th>Role</th>
            <th>Description</th>
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
import {getCache} from "@/assets/cache";
import {checkLogin, initModal, resetParsley} from '@/assets/forms.js'

export default {
  components: {
    confirmationModal,
    dataTable,
  },
  props: [],
  data() {
    return {
      isAdmin: getCache('session', 'admin') === 'true',
      usersEntries: [],
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
      selectedUserEntries: [],
      repoTableColumnDefs: [
        {"width": "10%", "targets": 0, "className": "text-center",},
        {"width": "10%", "targets": 1, "className": "text-center",},
        {"width": "20%", "targets": 2, "className": "text-center",},
        {"width": "15%", "targets": 3, "className": "text-center",},
        {"width": "30%", "targets": 4, "className": "text-left",},
        {"width": "15%", "targets": 5, "className": "text-center",},
      ],
    }
  },
  computed: {
    $form() {
      return $(this.$el).find('form');
    },
  },
  methods: {
    reloadPage() {
      this.getItems();
    },
    cellsData(e) {
      let link = `<a href="#" class="btn btn-primary btn-sm" @click="edit">
                                <span class="fas fa-edit"></span>Edit
                            </a>`;
      let link2 = `<a href="#" class="btn btn-primary btn-sm" @click="delete" data-toggle="modal"
                    data-refDate='${JSON.stringify(e)}' data-target="#confirm-deletion">
                                <span class="fas fa-trash"></span>Delete
                            </a>`;
      return [
        {'data-order': e.username, 'html': e.username},
        {'data-order': e.name, 'html': e.name},
        {'data-order': e.email, 'html': e.email},
        {'data-order': e.role, 'html': e.role},
        {'data-order': e.description, 'html': e.description},
        this.isAdmin ? link + link2 : link,
      ];
    },
    clickRow(e, $currentRow) {
      const target = $(e.target);
      const item = this.usersEntries[$currentRow];
      if (target.children().hasClass('fa-edit')) {
        this.$router.push({path: '/user/edit', query: {user: item.username}});
      } else if (target.children().hasClass('fa-trash')) {
        this.deleteItem(item);
      }
    },
    mountedData() {
      // debugger;
    },
    getItems() {
      let $this = this;
      let apiUrl = getCache('config', 'apiUrl');
      let url = apiUrl + '/users';
      if (!this.isAdmin) {
        const user = JSON.parse(getCache('session', 'user'));
        if (!user) return;
        url += '?username=' + user.username;
      }
      $('.loading').show();
      $.get(url).then(function (res) {
        $this.usersEntries = res;
        $('.loading').hide();
        $this.$refs.usersDataTable.mPromise.then(
            () => $this.$refs.usersDataTable.load()
        );
      });
    },
    deleteItem(item) {
      this.selectedUserEntries.push(item);
      this.$refs.confirmDeletion.show();
    },
    confirm(_this) {
      console.log(_this);
      let $this = this;
      $('.loading').show();
      let apiUrl = getCache('config', 'apiUrl');
      const url = apiUrl + '/user?username=' + this.selectedUserEntries[0].username;
      $.ajax({
        url: url,
        type: "DELETE",
        success: function () {
          $this.$refs.usersDataTable.dataTable.clear();
          $this.$refs.usersDataTable.dataTable.destroy();
          $this.getItems();
          $('.loading').hide();
          $this.$refs.confirmDeletion.hide();
        },
        error: function (e) {
          console.log(e);
        }
      });
    }
  },
  mounted() {
    checkLogin(this);
    this.getItems();
    resetParsley(this.$form);
    initModal();
  },
};

</script>