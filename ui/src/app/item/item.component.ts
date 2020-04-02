import {Component, OnInit, ViewChild} from '@angular/core';
import {ItemCreateComponent} from './item-create/item-create.component';
import {ItemListComponent} from './item-list/item-list.component';
import {SessionService} from '../shared/session.service';
import {ItemEditComponent} from './item-edit/item-edit.component';

@Component({
  selector: 'app-item',
  templateUrl: './item.component.html',
  styleUrls: ['./item.component.css']
})
export class ItemComponent implements OnInit {

  @ViewChild(ItemCreateComponent, {static: true})
  creationItem: ItemCreateComponent;

  @ViewChild(ItemListComponent, {static: true})
  listItem: ItemListComponent;

  @ViewChild(ItemEditComponent, {static: true})
  editItem: ItemEditComponent;

  permission;

  constructor(private sessionService: SessionService) {
  }

  ngOnInit() {
  }

  openModal(): void {
    this.creationItem.newItem();
  }

  createItem(created: boolean) {
    if (created) {
      this.listItem.listItem();
    }
  }

  updateItem(item) {
    this.editItem.editItem(item);
  }
}
