import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import * as globals from '../../globals';
import {Item} from '../item';
import {ItemService} from '../item.service';
import {AlertLevels} from '../../base/header/components/common-alert/alert';
import {CommonAlertService} from '../../base/header/common-alert.service';

@Component({
  selector: 'app-item-edit',
  templateUrl: './item-edit.component.html',
  styleUrls: ['./item-edit.component.css']
})
export class ItemEditComponent implements OnInit {


  @Output() edit = new EventEmitter<boolean>();
  editItemOpened = false;
  item = new Item();
  isSubmitGoing = false;
  name_pattern = globals.chinese_name_pattern;
  name_pattern_tip = globals.chinese_name_pattern_tip;
  oldName: string;

  constructor(private itemService: ItemService, private alert: CommonAlertService) {
  }

  ngOnInit() {
  }


  editItem(item) {
    this.item = item;
    this.oldName = this.item.name;
    this.editItemOpened = true;
  }

  onCancel() {
    this.editItemOpened = false;
  }

  onSubmit() {
    this.itemService.updateItem(this.item, this.oldName).subscribe(res => {
      this.alert.showAlert('更新成功', AlertLevels.SUCCESS);
      this.editItemOpened = false;
    });
  }
}
