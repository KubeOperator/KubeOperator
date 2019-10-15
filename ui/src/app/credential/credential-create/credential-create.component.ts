import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {Credential} from '../credential-list/credential';
import {CredentialService} from '../credential.service';
import {NgForm} from '@angular/forms';
import {CommonAlertService} from '../../base/header/common-alert.service';
import {AlertLevels} from '../../base/header/components/common-alert/alert';

@Component({
  selector: 'app-credential-create',
  templateUrl: './credential-create.component.html',
  styleUrls: ['./credential-create.component.css']
})
export class CredentialCreateComponent implements OnInit {

  @Output() create = new EventEmitter<boolean>();
  staticBackdrop = true;
  closable = false;
  createOpened: boolean;
  isSubmitGoing = false;
  item: Credential = new Credential();
  loading = false;
  @ViewChild('credentialForm', {static: true}) credentialForm: NgForm;

  constructor(private credentialService: CredentialService, private alertService: CommonAlertService) {
  }

  ngOnInit() {
  }

  onCancel() {
    this.createOpened = false;
    this.credentialForm.resetForm();
  }

  onSubmit() {
    if (this.isSubmitGoing) {
      return;
    }
    this.isSubmitGoing = true;
    this.loading = true;
    this.credentialService.createCredential(this.item).subscribe(data => {
      this.createOpened = false;
      this.isSubmitGoing = false;
      this.create.emit(true);
      this.loading = false;
      this.alertService.showAlert('创建凭据成功', AlertLevels.SUCCESS);
    }, err => {
      this.createOpened = false;
      this.isSubmitGoing = false;
      this.create.emit(true);
      this.loading = false;
      this.alertService.showAlert('创建凭据失败:' + err.reason + ' state code:' + err.status, AlertLevels.ERROR);
    });
  }

  newItem() {
    this.item = new Credential();
    this.createOpened = true;
  }


}
