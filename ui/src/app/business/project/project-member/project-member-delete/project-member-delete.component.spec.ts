import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ProjectMemberDeleteComponent } from './project-member-delete.component';

describe('ProjectMemberDeleteComponent', () => {
  let component: ProjectMemberDeleteComponent;
  let fixture: ComponentFixture<ProjectMemberDeleteComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ProjectMemberDeleteComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ProjectMemberDeleteComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
