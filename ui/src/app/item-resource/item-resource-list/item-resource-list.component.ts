import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {ItemResourceService} from '../item-resource.service';
import {ActivatedRoute} from '@angular/router';

@Component({
  selector: 'app-item-resource-list',
  templateUrl: './item-resource-list.component.html',
  styleUrls: ['./item-resource-list.component.css']
})
export class ItemResourceListComponent implements OnInit {

  loading = false;
  @Output() add = new EventEmitter();
  itemName;
  itemId;
  itemResources;

  constructor(private itemResourceService: ItemResourceService, private route: ActivatedRoute) {
  }

  ngOnInit() {
    this.itemName = this.route.snapshot.queryParams['name'];
    this.itemId = this.route.snapshot.queryParams['id'];
    this.getItemResources();
  }

  createItemResource(resourceType) {
    this.add.emit(resourceType);
  }

  getItemResources() {
    this.itemResourceService.getItemResources(this.itemName).subscribe(res => {
      this.itemResources = res;
    });
  }
}
