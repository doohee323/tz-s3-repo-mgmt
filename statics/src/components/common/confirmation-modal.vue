<template>
  <div class="modal fade" :id="id" :tabindex="tabIndex" role="dialog" :aria-labelledby="headerId">
    <div class="modal-dialog" role="document">
      <div class="modal-content">
        <form method="post">
          <input v-for="value in entries.map(paramValue)" type="hidden" :name="paramName" :value="value" :key="value">
          <div class="modal-header">
            <h5 class="modal-title" :id="headerId">{{ title }}</h5>
            <button type="button" class="close" data-dismiss="modal" :aria-label="messages.cancel">
              <span aria-hidden="true">&times;</span>
            </button>
          </div>
          <div class="modal-body text-center">
            <div v-show="isOneEntry">
              <slot name="one" :entry="entry"></slot>
            </div>
            <div v-show="!isOneEntry">
              <slot name="multiple" :entries="entries" :entry-count="entryCount"></slot>
            </div>
            <p v-if="permanent">This action cannot be undone.</p>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" data-dismiss="modal">
              Cancel
            </button>
            <button type="button" class="btn once" :class="`btn-${context}`" @click="confirm(this)">
              <slot name="confirm">Confirm</slot>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>

import $ from 'jquery';

export default {
  name: 'confirmation-modal',
  props: {
    id: {
      type: String,
      default: `confirm_1`,
    },
    title: {
      type: String,
      required: true,
    },
    entries: {
      type: Array,
      required: true,
    },
    context: {
      type: String,
      default: 'danger',
    },
    paramName: {
      type: String,
      required: true,
    },
    paramValue: {
      type: Function,
      default: e => e,
    },
    permanent: {
      type: Boolean,
      default: true,
    },
    tabIndex: {
      type: Number,
      default: -1,
    },
    confirm: {
      type: Function,
    },
  },
  data() {
    return {
      messages: {
        cancel: 'Cancel',
      },
    }
  },
  computed: {
    entry() {
      return this.entries[0] || {};
    },
    entryCount() {
      return this.entries.length;
    },
    headerId() {
      return `${this.id}_header`;
    },
    isOneEntry() {
      return this.entryCount === 1;
    },
  },
  methods: {
    show() {
      $(`#${this.id}`).modal('show');
    },
    hide() {
      $(`#${this.id}`).modal('hide');
    },
    data() {
      $(`#${this.id}`).modal('hide');
    },
  },
};

</script>