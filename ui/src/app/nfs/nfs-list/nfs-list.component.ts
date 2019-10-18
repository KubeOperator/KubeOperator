import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {NfsStorage} from '../nfs';
import {NfsService} from '../nfs.service';
import {fadeSlide} from '@clr/angular';
import {AlertLevels} from '../../base/header/components/common-alert/alert';
import {CommonAlertService} from '../../base/header/common-alert.service';

@Component({
  selector: 'app-nfs-list',
  templateUrl: './nfs-list.component.html',
  styleUrls: ['./nfs-list.component.css']
})
export class NfsListComponent implements OnInit {

  constructor(private nfsService: NfsService, private alertService: CommonAlertService) {
  }

  loading = false;
  selected: NfsStorage[] = [];
  items: NfsStorage[] = [];
  showDelete = false;
  resourceTypeName = 'NFS';

  @Output() add = new EventEmitter();

  ngOnInit() {
    this.listItems();
  }

  listItems() {
    this.nfsService.list().subscribe(data => {
      this.items = data;
    });
  }

  addNew() {
    this.add.emit();
  }

  refresh() {
    this.listItems();
  }

  delete() {
    const promises: Promise<{}>[] = [];
    this.selected.forEach(item => {
        promises.push(this.nfsService.delete(item.name).toPromise());
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

}
