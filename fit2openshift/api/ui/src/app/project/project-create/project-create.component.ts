import { Component, OnInit, ViewChild, OnDestroy, EventEmitter, Output } from '@angular/core';
import { NgForm } from '@angular/forms';
import { Subject } from 'rxjs';
import { debounceTime } from 'rxjs/operators';

import { Project } from '../project';
import { ProjectService } from '../project.service';

@Component({
  selector: 'app-project-create',
  templateUrl: './project-create.component.html',
  styles: []
})
export class ProjectCreateComponent implements OnInit, OnDestroy {
  @ViewChild('projectForm')
  currentForm: NgForm;
  createProjectOpened = false;
  isNameValid = true;
  staticBackdrop = true;
  closable = false;
  checkOnGoing = false;
  isSubmitOnGoing = false;
  projectForm: NgForm;
  hasChanged = false;
  nameTooltipText: string;
  project: Project = new Project();
  proNameChecker: Subject<string> = new Subject<string>();
  @Output() created = new EventEmitter<boolean>();

  constructor(private projectService: ProjectService) { }

  ngOnInit() {
    this.proNameChecker.pipe(
      debounceTime(300))
      .subscribe((name: string) => {
        const projectNameField = this.currentForm.controls['project_name'];
        if (projectNameField) {
          this.isNameValid = projectNameField.valid;
          if (this.isNameValid) {
            // Check exiting from backend
            this.projectService
              .checkProjectExists(projectNameField.value).toPromise()
              .then((exists) => {
                if (exists) {
                  this.isNameValid = false;
                  this.nameTooltipText = 'PROJECT.NAME_ALREADY_EXISTS';
                } else {
                  this.isNameValid = true;
                }
                this.checkOnGoing = false;
              })
              .catch(error => {
                this.checkOnGoing = false;
              });
          } else {
            this.nameTooltipText = 'PROJECT.NAME_TOOLTIP';
          }
        }
      });
  }

  onCancel() {
    this.createProjectOpened = false;
  }

  ngOnDestroy(): void {
    this.proNameChecker.unsubscribe();
  }

  public get isValid(): boolean {
    return this.currentForm &&
      this.currentForm.valid &&
      !this.isSubmitOnGoing &&
      this.isNameValid &&
      !this.checkOnGoing;
  }

  newProject() {
    this.project = new Project();
    this.hasChanged = false;
    this.isNameValid = true;
    this.createProjectOpened = true;
  }

  onSubmit() {
    if (this.isSubmitOnGoing) {
      return ;
    }

    this.isSubmitOnGoing = true;
    this.projectService
      .createProject(this.project)
      .subscribe(
        project => {
          this.isSubmitOnGoing = false;

          this.created.emit(true);
          // this.messageHandlerService.showSuccess('PROJECT.CREATED_SUCCESS');
          this.createProjectOpened = false;
        },
        error => {
          this.isSubmitOnGoing = false;

          // let errorMessage: string;
          // if (error instanceof Response) {
          //   switch (error.status) {
          //     case 409:
          //       this.translateService.get("PROJECT.NAME_ALREADY_EXISTS").subscribe(res => errorMessage = res);
          //       break;
          //     case 400:
          //       this.translateService.get("PROJECT.NAME_IS_ILLEGAL").subscribe(res => errorMessage = res);
          //       break;
          //     default:
          //       this.translateService.get("PROJECT.UNKNOWN_ERROR").subscribe(res => errorMessage = res);
          //   }
          //   this.messageHandlerService.handleError(error);
          // }
        });
  }

  handleValidation(): void {
    const projectNameField = this.currentForm.controls['project_name'];
    if (projectNameField) {
      this.proNameChecker.next(projectNameField.value);
    }
  }

}
