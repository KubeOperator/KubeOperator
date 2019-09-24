import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {PackageService} from '../../package/package.service';
import {Package} from '../../package/package';

@Component({
  selector: 'app-upgrade',
  templateUrl: './upgrade.component.html',
  styleUrls: ['./upgrade.component.css']
})
export class UpgradeComponent implements OnInit {
  opened = false;
  currentPackageName: string;
  currentPackage: Package;
  newPackage: Package;
  packages: Package[] = [];
  @Output() openedChange = new EventEmitter();
  @Output() paramsChange = new EventEmitter();
  @Output() confirm = new EventEmitter();

  constructor(private packageService: PackageService) {
  }

  ngOnInit() {
  }

  listPackage() {
    this.packageService.listPackage().subscribe(data => {
      this.packages = data;
      this.packages.forEach(p => {
        if (p.name === this.currentPackageName) {
          this.currentPackage = p;
        }
      });
    });
  }

  close() {
    this.opened = false;
    this.openedChange.emit(this.opened);
  }

  onConfirm() {
    this.paramsChange.emit();
    this.confirm.emit();
  }


}
