import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MultiClusterRepositoryErrorDetailComponent } from './multi-cluster-repository-error-detail.component';

describe('MultiClusterRepositoryErrorDetailComponent', () => {
  let component: MultiClusterRepositoryErrorDetailComponent;
  let fixture: ComponentFixture<MultiClusterRepositoryErrorDetailComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ MultiClusterRepositoryErrorDetailComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(MultiClusterRepositoryErrorDetailComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
