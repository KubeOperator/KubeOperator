import {Component, OnInit, ViewChild} from '@angular/core';
import {ItemResourceCreateComponent} from './item-resource-create/item-resource-create.component';
import {ItemResourceListComponent} from './item-resource-list/item-resource-list.component';

@Component({
  selector: 'app-item-resource',
  templateUrl: './item-resource.component.html',
  styleUrls: ['./item-resource.component.css']
})
export class ItemResourceComponent implements OnInit {

  @ViewChild(ItemResourceCreateComponent, {static: true})
  creation: ItemResourceCreateComponent;

  @ViewChild(ItemResourceListComponent, {static: true})
  listItemResource: ItemResourceListComponent;


  constructor() {
  }

  ngOnInit() {
  }

  openModal(event) {
    this.creation.createItemResource(event);
  }

}
