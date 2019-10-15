import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {RegionService} from '../region.service';
import {Region} from '../region';
import {RegionDetailComponent} from '../region-detail/region-detail.component';
import {CommonAlertService} from '../../base/header/common-alert.service';
import {AlertLevels} from '../../base/header/components/common-alert/alert';

@Component({
  selector: 'app-region-list',
  templateUrl: './region-list.component.html',
  styleUrls: ['./region-list.component.css']
})
export class RegionListComponent implements OnInit {

  items: Region[] = [];
  selected: Region[] = [];
  loading = true;
  showDelete = false;
  showDetail = false;
  resourceTypeName: '区域';
  @Output() add = new EventEmitter();
  @ViewChild(RegionDetailComponent, {static: true})
  child: RegionDetailComponent;

  constructor(private regionService: RegionService, private alertService: CommonAlertService) {
  }


  ngOnInit() {
    this.listItems();
  }

  listItems() {
    this.regionService.listRegion().subscribe(data => {
      this.items = data;
      this.loading = false;
    });
  }

  onShowDetail(item: Region) {
    this.showDetail = true;
    this.child.currentRegion = item;
  }

  delete() {
    const promises: Promise<{}>[] = [];
    this.selected.forEach(item => {
        promises.push(this.regionService.deleteRegion(item.name).toPromise());
      }
    );
    Promise.all(promises).then(data => {
      this.alertService.showAlert('删除成功', AlertLevels.SUCCESS);
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
