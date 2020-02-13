import {Component, OnInit} from '@angular/core';
import {ItemResourceService} from '../item-resource.service';
import {ActivatedRoute} from '@angular/router';

@Component({
  selector: 'app-item-resource-create',
  templateUrl: './item-resource-create.component.html',
  styleUrls: ['./item-resource-create.component.css']
})
export class ItemResourceCreateComponent implements OnInit {

  createOpened = false;
  itemName;
  resources = [];
  isSubmitGoing = false;


  constructor(private itemResourceService: ItemResourceService, private route: ActivatedRoute) {
  }

  ngOnInit() {
    this.itemName = this.route.snapshot.queryParams['name'];
  }

  createItemResource(resourceType) {
    this.itemResourceService.getResources(this.itemName, resourceType).subscribe(res => {
      this.resources = res;
      this.createOpened = true;
    });
  }

  onCancel() {
    this.createOpened = false;
  }

  onSubmit() {
    this.isSubmitGoing = true;
  }
}
