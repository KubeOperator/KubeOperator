import {Component, OnInit, ViewChild} from '@angular/core';
import {Package} from '../package';
import {PackageService} from '../package.service';
import {PackageLogoService} from '../package-logo.service';
import {PackageDetailComponent} from '../package-detail/package-detail.component';
import {CommonAlertService} from '../../base/header/common-alert.service';
import {AlertLevels} from '../../base/header/components/common-alert/alert';

@Component({
  selector: 'app-offline-list',
  templateUrl: './package-list.component.html',
  styleUrls: ['./package-list.component.css']
})
export class PackageListComponent implements OnInit {

  loading = true;
  packages: Package[] = [];
  selectedRow: Package[] = [];
  showDetail = false;
  @ViewChild(PackageDetailComponent, {static: true})
  child: PackageDetailComponent;

  constructor(private offlineService: PackageService, private packageLogoService: PackageLogoService, private alert: CommonAlertService) {
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

  onShowDetail(item) {
    this.showDetail = true;
    this.child.loadPackage(item);
  }

  getLogo(p: Package): string {
    return this.packageLogoService.getLogo(p.meta.resource);
  }


  refresh() {
    this.alert.showAlert('刷新成功', AlertLevels.SUCCESS);
    this.listOfflines();
  }

}
