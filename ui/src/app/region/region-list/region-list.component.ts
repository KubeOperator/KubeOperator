import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {Credential} from '../../credential/credential-list/credential';
import {RegionService} from '../region.service';
import {Region} from '../region';
import {TipLevels} from '../../tip/tipLevels';
import {TipService} from '../../tip/tip.service';
import {PackageDetailComponent} from '../../package/package-detail/package-detail.component';
import {RegionDetailComponent} from '../region-detail/region-detail.component';

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
  @ViewChild(RegionDetailComponent)
  child: RegionDetailComponent;

  constructor(private regionService: RegionService, private tipService: TipService) {
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
