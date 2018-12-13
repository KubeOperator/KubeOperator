import { Component, OnInit, ViewChild } from '@angular/core';
import { NgForm } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { Location } from '@angular/common';

import { Playbook } from '../playbook';
import { PlaybookService } from '../playbook.service';

@Component({
  selector: 'app-playbook-create',
  templateUrl: './playbook-create.component.html',
  styles: []
})
export class PlaybookCreateComponent implements OnInit {
  @ViewChild('playbookFrom') playbookForm: NgForm;
  playbook = new Playbook();
  projectName: string;

  constructor(
    private service: PlaybookService,
    private router: Router,
    private location: Location,
    private route: ActivatedRoute
  ) {
    this.projectName = this.route.snapshot.parent.params['project'];
  }

  ngOnInit() {
  }

  public get isValid(): boolean {
    return this.playbookForm &&
      this.playbookForm.valid;
  }

  onSubmit() {
    this.service.createPlaybook(this.playbook, this.projectName).subscribe(
      (playbook) => {
        console.log('Playbook created');
        this.location.back();
      },
      (error) => {
        if (error.status === 400) {
          console.log(this.playbookForm);

          for(const i in error.error) {
            console.log(i, error.error[i]);
          }
        }
        console.log(error);
      }
    );
  }
}
