import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ProjectMemberListComponent } from './project-member-list.component';

describe('ProjectMemberListComponent', () => {
  let component: ProjectMemberListComponent;
  let fixture: ComponentFixture<ProjectMemberListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ProjectMemberListComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ProjectMemberListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
