import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {UploadOutput, UploadInput, UploadFile, humanizeBytes, UploaderOptions} from 'ngx-uploader';
import {findLast} from '@angular/compiler/src/directive_resolver';
import {stringify} from '@angular/compiler/src/util';

@Component({
  selector: 'app-upload',
  templateUrl: './upload.component.html',
  styleUrls: ['./upload.component.css']
})
export class UploadComponent implements OnInit {
  options: UploaderOptions;
  files: UploadFile[];
  uploadInput: EventEmitter<UploadInput>;
  humanizeBytes: Function;
  // result
  @Output() uploaded: EventEmitter<any>;

  ngOnInit(): void {
  }

  constructor() {
    this.options = {concurrency: 5, maxUploads: 10, maxFileSize: 1000000};
    this.files = [];
    this.uploadInput = new EventEmitter<UploadInput>();
    this.uploaded = new EventEmitter<string[]>();
    this.humanizeBytes = humanizeBytes;
  }

  onUploadOutput(output: UploadOutput): void {
    switch (output.type) {
      case 'addedToQueue':
        if (typeof output.file !== 'undefined') {
          this.files.push(output.file);
        }
        break;
      case 'uploading':
        if (typeof output.file !== 'undefined') {
          // update current data in files array for uploading file
          const index = this.files.findIndex((file) => typeof output.file !== 'undefined' && file.id === output.file.id);
          this.files[index] = output.file;
        }
        break;
      case 'removed':
        this.files = this.files.filter((file: UploadFile) => file !== output.file);
        break;
      case 'done':
        let allCompleted = true;
        for (const file of this.files) {
          if (file.progress.status.valueOf() !== 2) {
            allCompleted = false;
            break;
          }
        }
        if (allCompleted) {
          const names = [];
          this.files.forEach(f => {
            names.push(f.name);
          });
          this.uploaded.emit(names);
        }
    }
  }

  startUpload(): void {
    const event: UploadInput = {
      type: 'uploadAll',
      url: '/api/v1/file/upload/',
      method: 'POST',
      data: {foo: 'bar'}
    };
    this.uploadInput.emit(event);
  }


  removeFile(id: string): void {
    for (const file of this.files) {
      if (file.id === id && file.progress.status.valueOf() !== 0) {
        return;
      }
    }
    this.uploadInput.emit({type: 'remove', id: id});
  }

  removeAllFiles(): void {
    this.uploadInput.emit({type: 'removeAll'});
  }
}
