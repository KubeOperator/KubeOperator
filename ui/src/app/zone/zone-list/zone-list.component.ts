import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {RegionService} from '../../region/region.service';
import {Zone} from '../zone';
import {ZoneService} from '../zone.service';
import {ZoneDetailComponent} from '../zone-detail/zone-detail.component';
import {CommonAlertService} from '../../base/header/common-alert.service';
import {AlertLevels} from '../../base/header/components/common-alert/alert';
import {ZoneEditComponent} from '../zone-edit/zone-edit.component';

@Component({
  selector: 'app-zone-list',
  templateUrl: './zone-list.component.html',
  styleUrls: ['./zone-list.component.css']
})
export class ZoneListComponent implements OnInit {

  items: Zone[] = [];
  selected: Zone[] = [];
  loading = false;
  showDelete = false;
  showEdit = false;
  showDetail = false;
  resourceTypeName: '可用区';
  @Output() add = new EventEmitter();
  @ViewChild(ZoneDetailComponent, {static: true}) child: ZoneDetailComponent;
  @ViewChild(ZoneEditComponent, {static: true}) childEdit: ZoneEditComponent;


  constructor(private regionService: RegionService, private zoneService: ZoneService, private alertService: CommonAlertService) {
  }


  ngOnInit() {
    this.listItems();
  }

  listItems() {
    this.zoneService.listZones().subscribe((data) => {
      this.items = data;
    });
  }

  delete() {
    const promises: Promise<{}>[] = [];
    this.selected.forEach(item => {
        promises.push(this.zoneService.deleteZone(item.name).toPromise());
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

  onShowDetail(item: Zone) {
    this.child.currentZone = item;
    this.showDetail = true;
  }

  onEdit(item: Zone) {
    this.zoneService.getZone(item.name).subscribe(data => {
      this.childEdit.item = data;
      this.childEdit.networkErrors = [];
      this.showEdit = true;
    });
  }

  refresh() {
    this.listItems();
  }

  addItem() {
    this.add.emit();
  }

}
