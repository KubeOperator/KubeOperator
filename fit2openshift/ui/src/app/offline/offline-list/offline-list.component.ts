import {Component, OnInit} from '@angular/core';
import {Offline} from '../Offline';
import {OfflineService} from '../offline.service';
import {MessageLevels} from '../../base/message/message-level';
import {TipService} from '../../tip/tip.service';
import {TipLevels} from '../../tip/tipLevels';

@Component({
  selector: 'app-offline-list',
  templateUrl: './offline-list.component.html',
  styleUrls: ['./offline-list.component.css']
})
export class OfflineListComponent implements OnInit {

  loading = true;
  offlines: Offline[] = [];
  selectedRow: Offline[] = [];

  constructor(private offlineService: OfflineService, private tipService: TipService) {
  }

  ngOnInit() {
    this.listOfflines();
  }

  listOfflines() {
    this.loading = true;
    this.offlineService.listOfflines().subscribe(data => {
      this.offlines = data;
      this.loading = false;
    });
  }

  refresh() {
    this.tipService.showTip('刷新成功', TipLevels.SUCCESS);
    this.listOfflines();
  }

}
