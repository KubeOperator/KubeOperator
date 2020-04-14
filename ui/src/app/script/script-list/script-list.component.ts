import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {Script} from '../script';
import {ScriptService} from '../script.service';
import {AlertLevels} from '../../base/header/components/common-alert/alert';
import {CommonAlertService} from '../../base/header/common-alert.service';

@Component({
  selector: 'app-script-list',
  templateUrl: './script-list.component.html',
  styleUrls: ['./script-list.component.css']
})
export class ScriptListComponent implements OnInit {

  selected: Script[] = [];
  items: Script[] = [];
  loading = true;
  deleteOpened = false;
  page = 1;
  size = 10;
  total = 100;

  @Output() create = new EventEmitter();


  constructor(private service: ScriptService, private alert: CommonAlertService) {
  }

  ngOnInit() {
    this.list();
  }

  onCreate() {
    this.create.emit();
  }

  onDelete() {
    this.deleteOpened = true;
  }


  delete() {
    const promises: Promise<{}>[] = [];
    this.selected.forEach(item => {
        promises.push(this.service.delete(item.name).toPromise());
      }
    );
    Promise.all(promises).then(data => {
      this.alert.showAlert('删除成功', AlertLevels.SUCCESS);
    }).finally(
      () => {
        this.deleteOpened = false;
        this.list();
        this.selected = [];
      }
    );
  }

  list() {
    this.service.list(this.page, this.size).subscribe(data => {
      this.items = data.results;
      this.total = data.count;
      this.loading = false;
    });
  }


}
