import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {CephService} from '../ceph.service';
import {AlertLevels} from '../../base/header/components/common-alert/alert';
import {CommonAlertService} from '../../base/header/common-alert.service';
import {Ceph} from '../ceph';

@Component({
  selector: 'app-ceph-list',
  templateUrl: './ceph-list.component.html',
  styleUrls: ['./ceph-list.component.css']
})
export class CephListComponent implements OnInit {

  items: Ceph[] = [];
  loading = false;
  resourceTypeName = 'Ceph';
  showDelete = false;
  selected: Ceph[] = [];

  @Output() add = new EventEmitter();


  constructor(private cephService: CephService, private alertService: CommonAlertService) {
  }

  ngOnInit() {
    this.listItems();
  }

  refresh() {
    this.listItems();
  }

  delete() {
    const promises: Promise<{}>[] = [];
    this.selected.forEach(item => {
        promises.push(this.cephService.delete(item.name).toPromise());
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

  addNew() {
    this.add.emit();
  }

  listItems() {
    this.cephService.list().subscribe(data => {
      this.items = data;
    });
  }
}
