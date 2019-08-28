import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {Zone} from '../zone';
import {Region} from '../../region/region';

@Component({
  selector: 'app-zone-detail',
  templateUrl: './zone-detail.component.html',
  styleUrls: ['./zone-detail.component.css']
})
export class ZoneDetailComponent implements OnInit {

  currentZone: Zone;
  @Input() showInfoModal = false;
  @Output() showInfoModalChange = new EventEmitter();

  constructor() {
  }

  ngOnInit() {
  }

  cancel() {
    this.showInfoModal = false;
    this.showInfoModalChange.emit(this.showInfoModal);
  }
}
