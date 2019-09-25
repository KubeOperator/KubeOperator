import {Component, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {PackageService} from '../../package/package.service';
import {Package} from '../../package/package';
import {NgForm} from '@angular/forms';

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
  @Output() confirm = new EventEmitter();
  @ViewChild('form', {static: true}) form: NgForm;

  constructor(private packageService: PackageService) {
  }

  reset() {
    this.newPackage = undefined;
    this.form.resetForm();
    this.listPackage();
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
      this.packages = this.packages.filter((p) => {
        return this.currentPackage.meta.vars['kube_version'] < p.meta.vars['kube_version'];
      });
    });
  }

  close() {
    this.opened = false;
    this.openedChange.emit(this.opened);
  }

  onConfirm() {
    this.confirm.emit();
  }


}
