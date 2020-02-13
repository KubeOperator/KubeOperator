import {Component, EventEmitter, OnInit, Output} from '@angular/core';

@Component({
  selector: 'app-item-resource-list',
  templateUrl: './item-resource-list.component.html',
  styleUrls: ['./item-resource-list.component.css']
})
export class ItemResourceListComponent implements OnInit {

  loading = false;
  @Output() add = new EventEmitter();


  constructor() {
  }

  ngOnInit() {
  }

  createItemResource(resourceType) {
    this.add.emit(resourceType);
  }
}
