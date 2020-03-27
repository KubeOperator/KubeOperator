import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {HostService} from '../host.service';
import {UploadComponent} from '../../shared/common-component/upload/upload.component';
import {HostCreateComponent} from '../host-create/host-create.component';

@Component({
  selector: 'app-host-import',
  templateUrl: './host-import.component.html',
  styleUrls: ['./host-import.component.css']
})
export class HostImportComponent implements OnInit {

  opened = false;
  file_names: string[] = [];

  @Output() imported: EventEmitter<boolean> = new EventEmitter();
  @ViewChild(UploadComponent, {static: true}) uploader: UploadComponent;

  constructor(private service: HostService) {
  }

  ngOnInit() {
  }

  open() {
    this.uploader.removeAllFiles();
    this.file_names = [];
    this.opened = true;
  }

  onUploaded(file_names: string[]) {
    this.file_names = file_names;
  }

  onCancel() {
    this.opened = false;
  }

  onSubmit() {
    this.service.import(this.file_names).subscribe(data => {
      this.imported.emit();
      this.opened = false;
    });
  }
}
