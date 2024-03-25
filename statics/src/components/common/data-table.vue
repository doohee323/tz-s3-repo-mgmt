<template>
  <div>
    <table class="table" :class="tableClasses" :name="tableId">
      <thead>
      <slot></slot>
      </thead>
      <tbody>
      </tbody>
    </table>
  </div>
</template>

<script>

import $ from 'jquery';

require('datatables.net-bs4');
require('datatables.net-bs4/css/dataTables.bootstrap4.css');
require('datatables.net-responsive-bs4');
require('datatables.net-responsive-bs4/css/responsive.bootstrap4.css');
require('datatables.net-select-bs4');
require('datatables.net-select-bs4/css/select.bootstrap4.min.css');

export default {
  name: 'data-table',
  data: function () {
    return {
      dataTable: null,
      tableId: this.$attrs.id,
    }
  },
  props: {
    entries: {
      type: Array,
      require: true,
    },
    cellsData: {
      type: Function,
      require: true,
    },
    actionButtons: {
      type: Array,
      default: null,
    },
    search: {
      type: Boolean,
      default: true,
    },
    info: {
      type: Boolean,
      default: true,
    },
    zeroRecords: {
      type: Boolean,
      default: true,
    },
    zeroRecordsMessage: {
      type: String,
      default: null,
    },
    pageLengths: {
      type: Array,
      default() {
        return [];
      },
    },
    order: {
      type: Array,
      default() {
        return [];
      },
    },
    ordering: {
      type: Boolean,
      default: true,
    },
    condensed: {
      type: Boolean,
      default: false,
    },
    responsive: {
      type: Boolean,
      default: true,
    },
    bordered: {
      type: Boolean,
      default: true,
    },
    striped: {
      type: Boolean,
      default: false,
    },
    classes: {
      type: String,
      default: null,
    },
    overflowWrap: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: null,
    },
    formInline: {
      type: Boolean,
      default: true,
    },
    autoWidth: {
      type: Boolean,
      default: true,
    },
    columnDefs: {
      type: Array,
      default() {
        return [];
      },
    },
    prohibit: {
      type: String,
      default: null,
    },
    minSelect: {
      type: Number,
      default: null,
    },
    maxSelect: {
      type: Number,
      default: null,
    },
    horizontalScroll: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    tableClasses() {
      let classDef = {
        'table-sm': this.condensed,
        'table-responsive-sm': this.responsive,
        'table-bordered': this.bordered,
        'table-striped': this.striped,
      };
      if (this.classes !== null) {
        for (let k of this.classes.split(' ')) {
          classDef[k] = true;
        }
      }
      return classDef;
    },
  },
  methods: {
    cellAttributes(d) {
      return typeof d === 'object' ?
          Object.keys(d).filter(k => k !== 'html').reduce((m, k) => Object.assign(m, {[k]: d[k]}), {}) : {};
    },
    cellMarkup(d) {
      return typeof d === 'object' ? d.html : d;
    },
    setSelectionNotes(dataTable, $selectionButtonGroup, $actionButton) {
      const buttonConfig = $actionButton.data('buttonConfig');
      const referencedCount = $(dataTable.rows({selected: true}).nodes()).find('.referenced').length;
      const permanentCount = $(dataTable.rows({selected: true}).nodes()).find('.permanent').length;
      let selectionNotes = [];
      if (referencedCount) {
        selectionNotes.push(referencedCount.toString());
      }
      if (permanentCount) {
        selectionNotes.push(permanentCount.toString());
      }
      const $selectionNotes = $selectionButtonGroup.parent().find('.selection-notes');
      $selectionNotes.html(selectionNotes.length ? '(' + selectionNotes.join(', ') + ')' : '');
      return $actionButton.text().trim() === buttonConfig.label && $selectionNotes.text() === '';
    },
    insertRow(e) {
      const $el = $(this.$el);
      const $tbody = $el.find('.dataTables_wrapper').length > 0 ? $el.find('.dataTables_wrapper').find('tbody') : $el.find('tbody');
      const $tr = $('<tr></tr>').appendTo($tbody).data('token', e.entry);
      for (const cd of e.cellsData(e.entry, 'insert', e.tableId).filter(cd => cd !== null)) {
        const tdEle = this.isNeedTop(cd) ? $('<td class="align-top"></td>') : $('<td></td>');
        tdEle.appendTo($tr).attr(this.cellAttributes(cd)).html(this.cellMarkup(cd));
      }
    },
    getDataElements(dataOnly) {
      this.dataElements = [];
      const rows = this.$table.find('tbody tr');
      for (let i = 0; i < rows.length; i++) {
        let dataElement = [];
        const rowData = $(rows[i]).find('td div');
        for (let j = 0; j < rowData.length; j++) {
          if ($(rowData[j]).children().length > 0) {
            const elem = $($(rowData[j]).children()[0]);
            if (dataOnly) {
              dataElement[elem.attr('name')] = $(elem).val();
            } else {
              dataElement[elem.attr('name')] = elem;
            }
          }
        }
        this.dataElements[i] = dataElement;
      }
      return this.dataElements;
    },
    getCurrentRowNum(e) {
      const rows = this.$table.find('tbody tr');
      const target = $(e.target);
      for (let i = 0; i < rows.length; i++) {
        if (target.closest('tbody')[0].rows[i] === target.closest('tr')[0]) {
          return i;
        }
      }
      return -1;
    },
    disabledColumn(currentRow, columnNames, disabled) {
      if (typeof columnNames === 'string') columnNames = columnNames.replace(/ /g, '').split(',');
      const data_fields = currentRow.find('.data-field');
      columnNames.map(s => {
        data_fields.map(i => {
          const field = $(data_fields[i]);
          if (field.find('input[name=' + s + ']').length > 0 ||
              field.find('select[name=' + s + ']').length > 0) {
            const elm = $(field.children()[0]);
            if (disabled) {
              elm.attr('disabled', true);
            } else {
              elm.removeAttr('disabled');
            }
          }
        });
      });
    },
    getColumn(currentRow, columnName) {
      const data_fields = currentRow.find('.data-field');
      return data_fields.map(i => {
        const field = $(data_fields[i]);
        if (field.find('input[name=' + columnName + ']').length > 0 ||
            field.find('select[name=' + columnName + ']').length > 0) {
          return field.children()[0];
        }
      });
    },
    resetColumns() {
      const elements = this.getDataElements();
      for (let i = 0; i < elements.length; i++) {
        for (let j in elements[i]) {
          elements[i][j].removeAttr('disabled').removeAttr('required');
        }
      }
    },
    clear() {
      this.dataTable.clear().draw();
    },
    draw(entries) {
      const $el = $(this.$el);
      const $table = $el.find('.table');
      const $tbody = $table.find('tbody');
      for (const entry of entries) {
        let data = [];
        const $tr = $(`<tr></tr>`).attr('data-prohibit', entry.$prohibit ? this.prohibit : '')
            .appendTo($tbody).data('token', entry);
        for (const cd of this.cellsData(entry).filter(cd => cd !== null)) {
          const tdEle = this.isNeedTop(cd) ? $('<td class="align-top"></td>') : $('<td></td>');
          const td = tdEle.appendTo($tr).attr(this.cellAttributes(cd)).html(this.cellMarkup(cd));
          data.push(td.html());
        }
        this.dataTable.row.add(data);
      }
      this.dataTable.draw();
    },
    isNeedTop(cd) {
      let html = cd.html ? cd.html : (cd && typeof cd === 'string' ? cd : '');
      return html.indexOf('btn-sm') > -1 || html.indexOf('error-handler') > -1;
    },
    load() {
      const language = {
        search: 'Search:',
      };
      if (this.zeroRecordsMessage !== null) {
        language.zeroRecords = this.zeroRecordsMessage;
      }
      const options = {
        dom: `<'row'<'col-md-6 selection-buttons'><'#manage_dom.col-md-6'lf>>` +
            `<'row'<'col-md-12'tr>>` +
            `<'row'<'col-md-6'i><'col-md-6'p>>`,
        language: language,
        paging: this.pageLengths !== null,
        select: this.actionButtons ? 'multi' : false,
        searching: this.search,
        ordering: this.order !== null,
        order: this.order,
        info: this.info,
        zeroRecords: this.zeroRecords,
        autoWidth: this.autoWidth,
        columnDefs: this.columnDefs,
      };

      if (this.pageLengths !== null && this.pageLengths.length > 0) {
        options.lengthMenu = this.pageLengths;
      }

      const $el = $(this.$el);
      const $table = $el.find('> table');
      const $tbody = $table.find('tbody');
      for (const entry of this.entries) {
        const $tr = $(`<tr></tr>`).attr('data-prohibit', entry.$prohibit ? this.prohibit : '')
            .appendTo($tbody).data('token', entry);
        for (const cd of this.cellsData(entry).filter(cd => cd !== null)) {
          const tdEle = this.isNeedTop(cd) ? $('<td class="align-top"></td>') : $('<td></td>');
          tdEle.appendTo($tr).attr(this.cellAttributes(cd)).html(this.cellMarkup(cd));
        }
      }

      $el.on('init.dt', $.proxy(function (e) {
        this.$emit('loadeddata', e, this);
      }, this));

      const dataTable = $table.DataTable(options)
          .on('user-select', (e, _dt, _type, _cell, originalEvent) => {
                let target = originalEvent.target, node = null;
                do {
                  node = target.nodeName.toLowerCase();
                  if (node === 'a') {
                    e.preventDefault();
                    return;
                  }
                  if ($(target).closest('td').find('.unselect').length > 0) {
                    e.preventDefault();
                    return;
                  }
                  target = $(target).parent()[0];
                } while (target && node !== 'td');
              }
          );

      if (this.overflowWrap) {
        $table.find('td').css('overflow-wrap', 'anywhere');
      }

      if (!this.formInline) {
        $table.closest('.dataTables_wrapper').removeClass('form-inline');
      }

      if (this.actionButtons) {

        const $selectionButtonGroup = $(`<div class='btn-group'></div>`)
            .appendTo($el.find('.dataTables_wrapper > .row > .selection-buttons'));

        $.each([
              ['Select All', 'select', 'check-square'],
              ['Deselect All', 'deselect', 'square']],
            (_i, [t, c, i]) => $selectionButtonGroup.append(
                `<button type="button" class="${c} btn btn-secondary" ${this.entries ? '' : 'disabled'}>` +
                `    <span class="fas fa-${i}"></span> ${t}` +
                `</button>`
            )
        );
        const $selectAllButton = $selectionButtonGroup.find('button.select');
        $selectAllButton.click(() => dataTable.rows(':visible').select());
        const $deselectAllButton = $selectionButtonGroup.find('button.deselect');
        $deselectAllButton.click(() => dataTable.rows(':visible').deselect());

        const $actionButtons = $(...this.actionButtons.map($.proxy(function (buttonConfig) {
          return $(`<a href="#" class="btn btn-${buttonConfig.context}"
                              data-min-select="${this.minSelect}" data-max-select="${this.maxSelect}"
                              data-target="#${this.prohibit}"></a>`)
              .data('buttonConfig', buttonConfig)
              .appendTo($selectionButtonGroup)
              .append(`<span class="fas fa-${buttonConfig.icon}"></span>`)
              .append(` ${buttonConfig.label}`)
              .click($.proxy(function (e) {
                if ($(e.target).attr('disabled') !== 'disabled' && !$(e.target).hasClass('disabled')) {
                  this.$emit(
                      buttonConfig.event,
                      $.map(dataTable.rows({selected: true}).nodes(), n => $(n).data('token')));
                }
              }, this));
        }, this)));

        const updateSelection = $.proxy(function () {
          const $wrapper = $el.find('.dataTables_wrapper');
          $selectAllButton.prop('disabled', !dataTable.rows(':visible', {selected: false}).count());
          $deselectAllButton.prop('disabled', !dataTable.rows(':visible', {selected: true}).count());

          const selectedEntries = $.map(dataTable.rows({selected: true}).nodes(), n => $(n).data('token'));
          const selectedCount = selectedEntries.length;

          let prohibited = {};
          $.each(dataTable.rows({selected: true}).nodes(), function (_, r) {
            let $selectedRow = $(r);
            let prohibitedValue = $selectedRow.data('prohibit');
            if (prohibitedValue) {
              $.each(prohibitedValue.split(' '), function (_, p) {
                prohibited[p] = true;
              });
            }
          });
          $wrapper.find('a[data-min-select], a[data-max-select]').each(function (_, link) {
            let $link = $(link);
            let operation = $link.data('target');
            operation = operation ? operation.substr(1) : $link.data('operation');
            let min = $link.data('minSelect') || 0;
            let max = $link.data('maxSelect');
            if (selectedCount < min || (max && selectedCount > max) || operation in prohibited) {
              $link.addClass('disabled');
            } else {
              $link.removeClass('disabled');
            }
          });
          $.each($actionButtons, (_i, actionButton) => {
            const $actionButton = $(actionButton);
            const buttonConfig = $actionButton.data('buttonConfig');
            const atLeast = buttonConfig.atLeast || 1;
            const atMost = buttonConfig.atMost;

            const isSelectionNotes = this.setSelectionNotes(dataTable, $selectionButtonGroup, $actionButton);
            const enable = (selectedCount >= atLeast) &&
                (typeof atMost === 'undefined' || selectedCount <= atMost) &&
                isSelectionNotes
            $actionButton.attr('disabled', !enable);
            $actionButton.addClass(enable ? null : 'disabled');
          });

          this.$emit('select', selectedEntries);

        }, this);
        updateSelection();
      }

      if (!this.zeroRecords) {
        $el.find('.dataTables_empty').closest('tr').remove();
      }

      $tbody.on('click', $.proxy(function (e) {
        this.$emit('click', e, this.getCurrentRowNum(e), this);
      }, this));
      $tbody.on('change keyup', $.proxy(function (e) {
        if ([9, 16, 29, 37, 39, 65, 67, 91].includes(e.keyCode)) return true;
        this.$emit('change', e, this.getCurrentRowNum(e), this);
      }, this));

      if (this.horizontalScroll) {
        $el.find('table').parent().css('overflow', 'auto');
      }

      this.dataTable = dataTable;
      this.$table = $table;
    },
  },
  mounted() {
    this.mPromise = new Promise((resolve) => resolve());
  },
};

</script>

