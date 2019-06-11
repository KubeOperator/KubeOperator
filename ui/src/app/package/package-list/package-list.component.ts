import {Component, OnInit} from '@angular/core';
import {Package} from '../package';
import {PackageService} from '../package.service';
import {MessageLevels} from '../../base/message/message-level';
import {TipService} from '../../tip/tip.service';
import {TipLevels} from '../../tip/tipLevels';
import {PackageLogoService} from '../package-logo.service';

@Component({
  selector: 'app-offline-list',
  templateUrl: './package-list.component.html',
  styleUrls: ['./package-list.component.css']
})
export class PackageListComponent implements OnInit {

  loading = true;
  packages: Package[] = [];
  selectedRow: Package[] = [];

  constructor(private offlineService: PackageService, private tipService: TipService, private packageLogoService: PackageLogoService) {
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

  getLogo(p: Package): string {
    return this.packageLogoService.getLogo(p.meta.resource);
  }


  refresh() {
    this.tipService.showTip('刷新成功', TipLevels.SUCCESS);
    this.listOfflines();
  }

}
