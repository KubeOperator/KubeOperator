import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {Package} from '../package';

@Component({
  selector: 'app-package-detail',
  templateUrl: './package-detail.component.html',
  styleUrls: ['./package-detail.component.css']
})
export class PackageDetailComponent implements OnInit {

  hostId: string;
  loading = false;
  currentPackage: Package = null;
  @Input() showInfoModal = false;
  @Output() showInfoModalChange = new EventEmitter();

  constructor() {
  }

  ngOnInit() {
  }

  loadPackage(pkg: Package) {
    console.log(pkg);
    setTimeout(() => {
      this.currentPackage = pkg;
      console.log(123);
      console.log(this.currentPackage);
    }, 10);
  }

  cancel() {
    this.showInfoModal = false;
    this.currentPackage = null;
    this.showInfoModalChange.emit(this.showInfoModal);
  }

}
