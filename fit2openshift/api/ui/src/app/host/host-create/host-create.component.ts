import { Component, OnInit, ViewChild, Output, EventEmitter } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { NgForm } from '@angular/forms';

import { HostService } from '../host.service';
import { Host } from '../host';

@Component({
  selector: 'app-host-create',
  templateUrl: './host-create.component.html',
  styles: []
})
export class HostCreateComponent implements OnInit {
  @ViewChild('hostForm')
  currentForm: NgForm;
  @Output() created = new EventEmitter<boolean>();
  host: Host = new Host();
  shown = false;
  projectName: string;

  constructor(private service: HostService, private route: ActivatedRoute) { }

  ngOnInit() {
    this.projectName = this.route.snapshot.parent.params['project'];
  }

  newHost(): void {
    this.host = new Host();
    this.host.port = 22;
    this.host.username = 'root';
    this.shown = true;
  }

  public get isValid(): boolean {
    return this.currentForm &&
      this.currentForm.valid;
  }

  onCancel() {
    this.shown = false;
  }

  onSubmit() {
    this.service.createHost(this.projectName, this.host).subscribe(
      host => {
        console.log(host);
        this.shown = false;
        this.created.emit(true);
      }
    );
  }

}
