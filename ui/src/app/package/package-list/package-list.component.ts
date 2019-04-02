import {Component, OnInit} from '@angular/core';
import {Package} from '../package';
import {PackageService} from '../package.service';
import {MessageLevels} from '../../base/message/message-level';
import {TipService} from '../../tip/tip.service';
import {TipLevels} from '../../tip/tipLevels';

@Component({
  selector: 'app-offline-list',
  templateUrl: './package-list.component.html',
  styleUrls: ['./package-list.component.css']
})
export class PackageListComponent implements OnInit {

  loading = true;
  packages: Package[] = [];
  selectedRow: Package[] = [];

  constructor(private offlineService: PackageService, private tipService: TipService) {
  }

  ngOnInit() {
    this.listOfflines();
  }

  listOfflines() {
    this.loading = true;
    this.offlineService.listPackage().subscribe(data => {
      this.packages = data;
      this.loading = false;
    });
  }

  refresh() {
    this.tipService.showTip('刷新成功', TipLevels.SUCCESS);
    this.listOfflines();
  }

}
