import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { StorageClassComponent } from './storage-class.component';

describe('StorageClassComponent', () => {
  let component: StorageClassComponent;
  let fixture: ComponentFixture<StorageClassComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ StorageClassComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(StorageClassComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
