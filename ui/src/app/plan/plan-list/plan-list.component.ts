import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {Region} from '../../region/region';
import {RegionService} from '../../region/region.service';
import {TipService} from '../../tip/tip.service';
import {TipLevels} from '../../tip/tipLevels';
import {Plan} from '../plan';
import {PlanService} from '../plan.service';
import {PlanDetailComponent} from '../plan-detail/plan-detail.component';

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
  showDetail = false;
  resourceTypeName: '计划';
  @Output() add = new EventEmitter();
  @ViewChild(PlanDetailComponent)
  child: PlanDetailComponent;

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

  onShowDetail(item: Plan) {
    this.showDetail = true;
    this.child.currentPlan = item;
  }

  refresh() {
    this.listItems();
  }

  addItem() {
    this.add.emit();
  }

  getDeployName(name: string) {
    switch (name) {
      case 'SINGLE':
        return '一主多节点';
      case 'MULTIPLE':
        return '多主多节点';
      default:
        return '无';
    }
  }
}
