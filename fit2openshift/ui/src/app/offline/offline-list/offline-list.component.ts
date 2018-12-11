import {Component, OnInit} from '@angular/core';
import {Offline} from '../Offline';
import {OfflineService} from '../offline.service';

@Component({
  selector: 'app-offline-list',
  templateUrl: './offline-list.component.html',
  styleUrls: ['./offline-list.component.css']
})
export class OfflineListComponent implements OnInit {

  loading = true;
  offlines: Offline[] = [];
  selectedRow: Offline[] = [];

  constructor(private offlineService: OfflineService) {
  }

  ngOnInit() {
    this.listOfflines();
  }

  listOfflines() {
    this.offlineService.listOfflines().subscribe(data => {
      this.offlines = data;
      this.loading = false;
    });
  }

}
