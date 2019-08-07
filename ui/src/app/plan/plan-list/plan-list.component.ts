import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {Region} from '../../region/region';
import {RegionService} from '../../region/region.service';
import {TipService} from '../../tip/tip.service';
import {TipLevels} from '../../tip/tipLevels';
import {Plan} from '../plan';
import {PlanService} from '../plan.service';

@Component({
  selector: 'app-plan-list',
  templateUrl: './plan-list.component.html',
  styleUrls: ['./plan-list.component.css']
})
export class PlanListComponent implements OnInit {

  items: Plan[] = [];
  selected: Plan[] = [];
  loading = true;
  showDelete = false;
  resourceTypeName: '计划';
  @Output() add = new EventEmitter();

  constructor(private regionService: RegionService, private tipService: TipService,
              private planService: PlanService) {
  }

  ngOnInit() {
    this.listItems();
  }

  listItems() {
    this.planService.listPlan().subscribe(data => {
      this.loading = false;
      this.items = data;
    });
  }

  delete() {
    const promises: Promise<{}>[] = [];
    this.selected.forEach(item => {
        promises.push(this.planService.deletePlan(item.name).toPromise());
      }
    );
    Promise.all(promises).then(data => {
      this.tipService.showTip('删除成功', TipLevels.SUCCESS);
    }, error => {
      this.tipService.showTip('删除失败' + error.toString(), TipLevels.ERROR);
    }).finally(
      () => {
        this.showDelete = false;
        this.listItems();
        this.selected = [];
      }
    );
  }

  refresh() {
    this.listItems();
  }

  addItem() {
    this.add.emit();
  }

}
