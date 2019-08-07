import {Component, OnInit, ViewChild} from '@angular/core';
import {RegionCreateComponent} from '../region/region-create/region-create.component';
import {RegionListComponent} from '../region/region-list/region-list.component';
import {PlanCreateComponent} from './plan-create/plan-create.component';
import {PlanListComponent} from './plan-list/plan-list.component';

@Component({
  selector: 'app-plan',
  templateUrl: './plan.component.html',
  styleUrls: ['./plan.component.css']
})
export class PlanComponent implements OnInit {

  @ViewChild(PlanCreateComponent)
  creation: PlanCreateComponent;

  @ViewChild(PlanListComponent)
  listRegion: PlanListComponent;

  constructor() {
  }

  ngOnInit() {
  }

  openModal() {
    this.creation.newItem();
  }

  create(created: boolean) {
    if (created) {
      this.refresh();
    }
  }

  refresh() {
    this.listRegion.refresh();
  }

}
