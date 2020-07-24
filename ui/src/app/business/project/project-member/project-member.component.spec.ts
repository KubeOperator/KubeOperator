import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ProjectMemberComponent } from './project-member.component';

describe('ProjectMemberComponent', () => {
  let component: ProjectMemberComponent;
  let fixture: ComponentFixture<ProjectMemberComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ProjectMemberComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ProjectMemberComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
