import {Component, OnInit} from '@angular/core';
import {Offline} from '../Offline';
import {OfflineService} from '../offline.service';
import {MessageService} from '../../base/message.service';
import {MessageLevels} from '../../base/message/message-level';

@Component({
  selector: 'app-offline-list',
  templateUrl: './offline-list.component.html',
  styleUrls: ['./offline-list.component.css']
})
export class OfflineListComponent implements OnInit {

  loading = true;
  offlines: Offline[] = [];
  selectedRow: Offline[] = [];

  constructor(private offlineService: OfflineService, private messageService: MessageService) {
  }

  ngOnInit() {
    this.messageService.announceMessage('text', MessageLevels.ERROR);
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
    this.listOfflines();
  }

}
