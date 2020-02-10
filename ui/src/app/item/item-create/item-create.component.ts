import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import * as globals from '../../globals';
import {Item} from '../item';
import {AlertLevels} from "../../base/header/components/common-alert/alert";
import {ItemService} from "../item.service";
import {CommonAlertService} from "../../base/header/common-alert.service";


@Component({
  selector: 'app-item-create',
  templateUrl: './item-create.component.html',
  styleUrls: ['./item-create.component.css']
})
export class ItemCreateComponent implements OnInit {


  constructor(private itemService: ItemService, private alert: CommonAlertService) {
  }

  createItemOpened: boolean;
  isSubmitGoing = false;
  item: Item = new Item();
  name_pattern = globals.cluster_name_pattern;
  name_pattern_tip = globals.cluster_name_pattern_tip;
  @Output() create = new EventEmitter<boolean>();


  ngOnInit() {

  }

  newItem() {
    this.createItemOpened = true;
  }

  onCancel() {
    this.createItemOpened = false;
  }

  onSubmit() {
    if (this.isSubmitGoing) {
      return;
    }
    this.isSubmitGoing = true;

    this.itemService.createItem(this.item).subscribe(data => {
      this.createItemOpened = false;
      this.isSubmitGoing = false;
      this.create.emit(true);
      this.alert.showAlert('创建项目成功', AlertLevels.SUCCESS);
    }, res => {
      this.createItemOpened = false;
      this.isSubmitGoing = false;
      this.create.emit(true);
    });
  }
}
